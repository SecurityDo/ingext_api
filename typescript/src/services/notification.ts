import type { IngextClient } from "../client.js";
import type {
  GenericDAORequest,
  GenericDaoAddResponse,
  GenericDaoEntryResponse,
  GenericDaoListResponse,
} from "../types/dao.js";
import type { EndpointConfig } from "../types/notification.js";
import type { FPLScript } from "../types/platform.js";
import { PlatformService } from "./platform.js";

const DS = "api/ds";

/** The `actionConfig.target` of every action a notification endpoint can use. */
const NOTIFICATION_TARGET = "Platform Notification";

/**
 * Wraps `platform_notification_endpoint_dao`, the dao that holds the
 * notification endpoints an FPL action delivers through.
 *
 * An endpoint pairs an `integration` ("Email", "Slack") with the name of an
 * fpl_action script in its `action` field, plus the per-integration config that
 * script reads as its `config` argument. `listActions` reports which action
 * names exist and which integration each one expects.
 */
export class NotificationService {
  constructor(private client: IngextClient) {}

  /**
   * Store a new endpoint. The dao keys endpoints by name and rejects a name it
   * already holds with "duplicate endpoint" — add is not an upsert; use
   * `update` for that.
   *
   * The returned id is empty in practice: the dao answers an add with `{}` and
   * no id field. Address the endpoint by its name, not by this value.
   */
  async add(entry: EndpointConfig): Promise<string> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "add",
      args: { entry },
    };
    const res = await this.client.call<GenericDaoAddResponse>(
      DS,
      "platform_notification_endpoint_dao",
      req,
    );
    // The dao answers "{}" today, so keep the declared string type honest
    // rather than handing back undefined.
    return res.id ?? "";
  }

  /**
   * Replace the stored endpoint of `entry.name` with `entry`.
   *
   * Two things to know. The dao takes the whole entry, not a patch: a field left
   * out is stored empty, so read the current endpoint with `get` and modify that
   * rather than building one from scratch. And update is an UPSERT — a name the
   * dao does not hold is created rather than refused, so a typo in `entry.name`
   * silently adds a second endpoint instead of editing the one meant.
   */
  async update(entry: EndpointConfig): Promise<void> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "update",
      args: { id: entry.name, entry },
    };
    await this.client.call(DS, "platform_notification_endpoint_dao", req);
  }

  /**
   * Add an Email endpoint. `action` names the fpl_action script that renders
   * the mail, e.g. "Generic_Email_Action". See `add` on the returned id.
   */
  async addEmail(
    name: string,
    action: string,
    to: string[],
    cc: string[],
  ): Promise<string> {
    return await this.add({
      name,
      integration: "Email",
      action,
      email: { to, cc },
    });
  }

  /**
   * Add a Slack endpoint. `integrationName` is the name of the configured Slack
   * integration (`platform_integration_dao`) holding the token to post with.
   */
  async addSlack(
    name: string,
    action: string,
    integrationName: string,
    channels: string[],
  ): Promise<string> {
    return await this.add({
      name,
      integration: "Slack",
      action,
      slack: {
        integrationName,
        channels,
        ...(channels.length === 1 ? { channel: channels[0] } : {}),
      },
    });
  }

  /**
   * One endpoint by name.
   *
   * A name the dao does not hold **throws** ("export not found: <name>") rather
   * than returning null. The null return is only for a dao that answers with a
   * null entry, which the server is not observed to do.
   */
  async get(name: string): Promise<EndpointConfig | null> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "get",
      args: { id: name },
    };
    const res = await this.client.call<GenericDaoEntryResponse<EndpointConfig>>(
      DS,
      "platform_notification_endpoint_dao",
      req,
    );
    return res.entry ?? null;
  }

  /**
   * Remove the endpoint of that name. Deleting a name the dao does not hold
   * throws ("unknown export") rather than succeeding quietly.
   */
  async delete(name: string): Promise<void> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, "platform_notification_endpoint_dao", req);
  }

  async list(): Promise<EndpointConfig[]> {
    const req: GenericDAORequest<EndpointConfig> = { action: "list" };
    const res = await this.client.call<GenericDaoListResponse<EndpointConfig>>(
      DS,
      "platform_notification_endpoint_dao",
      req,
    );
    return res.entries ?? [];
  }

  /**
   * The actions an endpoint can name: the fpl_action scripts whose
   * `actionConfig.target` is "Platform Notification". Pass `integration` to keep
   * only the actions expecting that endpoint kind ("Email", "Slack").
   */
  async listActions(integration?: string): Promise<FPLScript[]> {
    const actions = await new PlatformService(this.client).listActions();
    return actions.filter(
      (a) =>
        a.actionConfig?.target === NOTIFICATION_TARGET &&
        (!integration || a.actionConfig?.integration === integration),
    );
  }
}
