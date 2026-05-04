import type { IngextClient } from "../client.js";
import type { GenericDAORequest, GenericDaoListResponse } from "../types/dao.js";
import type { GithubRepo } from "../types/repo.js";

const DS = "api/ds";

/**
 * Minimal mirror of `github.com/google/go-github/v64/github.RepositoryContent`.
 * The Go server emits the GitHub API content shape unchanged. Only the most
 * commonly used fields are included; unknown fields are preserved as `unknown`
 * via index access.
 */
export interface RepositoryContent {
  type?: string;
  encoding?: string;
  size?: number;
  name?: string;
  path?: string;
  content?: string;
  sha?: string;
  url?: string;
  git_url?: string;
  html_url?: string;
  download_url?: string;
  [key: string]: unknown;
}

export interface GitRepoContentRequest {
  id?: string;
  path: string;
}

export interface GitRepoContentResponse {
  file: RepositoryContent | null;
  dirs: RepositoryContent[];
}

export interface ImportRepoObjectsRequest {
  id?: string;
  paths: string[];
}

export class RepoService {
  constructor(private client: IngextClient) {}

  async listRepos(): Promise<GithubRepo[]> {
    const req: GenericDAORequest<GithubRepo> = { action: "list" };
    const res = await this.client.call<GenericDaoListResponse<GithubRepo>>(
      DS,
      "github_repo_dao",
      req,
    );
    return res.entries ?? [];
  }

  async getRepoContent(repoId: string, repoPath: string): Promise<GitRepoContentResponse> {
    const req: GitRepoContentRequest = { id: repoId, path: repoPath };
    return await this.client.call<GitRepoContentResponse>(
      DS,
      "github_repo_get_content",
      req,
    );
  }

  async importRepoProcessors(repoId: string, filePaths: string[]): Promise<void> {
    const req: ImportRepoObjectsRequest = { id: repoId, paths: filePaths };
    await this.client.call(DS, "platform_import_processors", req);
  }

  async importAppTemplates(repoId: string, filePaths: string[]): Promise<void> {
    const req: ImportRepoObjectsRequest = { id: repoId, paths: filePaths };
    await this.client.call(DS, "platform_import_application_templates", req);
  }

  async importLakeSchemas(repoId: string, filePaths: string[]): Promise<void> {
    const req: ImportRepoObjectsRequest = { id: repoId, paths: filePaths };
    await this.client.call(DS, "platform_import_lake_schemas", req);
  }
}
