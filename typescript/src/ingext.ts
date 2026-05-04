import { IngextClient, type IngextClientOptions } from "./client.js";
import { ApplicationService } from "./services/application.js";
import { AuthService } from "./services/auth.js";
import { CollectorService } from "./services/collector.js";
import { DatalakeService } from "./services/datalake.js";
import { EventWatchService } from "./services/eventwatch.js";
import { FPLService } from "./services/fpl.js";
import { GridService } from "./services/grid.js";
import { NotificationService } from "./services/notification.js";
import { PlatformService } from "./services/platform.js";
import { RepoService } from "./services/repo.js";
import { ResourceService } from "./services/resource.js";
import { SearchService } from "./services/search.js";
import { SyslogService } from "./services/syslog.js";

/**
 * Top-level convenience client. Constructs a single `IngextClient` and
 * exposes one ready-to-use service per backend area. Mirrors the Go
 * pattern of `NewIngextClient` + `NewAuthService(client)` etc., but in a
 * single object for ergonomic TypeScript use.
 */
export class Ingext {
  readonly client: IngextClient;

  readonly auth: AuthService;
  readonly datalake: DatalakeService;
  readonly search: SearchService;
  readonly eventwatch: EventWatchService;
  readonly fpl: FPLService;
  readonly syslog: SyslogService;
  readonly grid: GridService;
  readonly collector: CollectorService;
  readonly resource: ResourceService;
  readonly notification: NotificationService;
  readonly application: ApplicationService;
  readonly repo: RepoService;
  readonly platform: PlatformService;

  constructor(opts: IngextClientOptions) {
    this.client = new IngextClient(opts);

    this.auth = new AuthService(this.client);
    this.datalake = new DatalakeService(this.client);
    this.search = new SearchService(this.client);
    this.eventwatch = new EventWatchService(this.client);
    this.fpl = new FPLService(this.client);
    this.syslog = new SyslogService(this.client);
    this.grid = new GridService(this.client);
    this.collector = new CollectorService(this.client);
    this.resource = new ResourceService(this.client);
    this.notification = new NotificationService(this.client);
    this.application = new ApplicationService(this.client);
    this.repo = new RepoService(this.client);
    this.platform = new PlatformService(this.client);
  }

  setToken(token: string): void {
    this.client.setToken(token);
  }

  setDebug(flag: boolean): void {
    this.client.setDebug(flag);
  }

  async close(): Promise<void> {
    await this.client.close();
  }
}
