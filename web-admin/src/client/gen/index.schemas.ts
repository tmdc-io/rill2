/**
 * Generated by orval v6.21.0 🍺
 * Do not edit manually.
 * rill/admin/v1/api.proto
 * OpenAPI spec version: version not set
 */
export type AdminServiceSearchUsersParams = {
  emailPattern?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListBookmarksParams = {
  projectId?: string;
};

export type AdminServiceGetUserParams = {
  email?: string;
};

export type AdminServiceSudoGetResourceParams = {
  userId?: string;
  orgId?: string;
  projectId?: string;
  deploymentId?: string;
  instanceId?: string;
};

export type AdminServiceSearchProjectNamesParams = {
  namePattern?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetReportMetaParams = {
  branch?: string;
  report?: string;
  /**
   * This is a request variable of the map type. The query format is "map_name[key]=value", e.g. If the map name is Age, the key type is string, and the value type is integer, the query parameter is expressed as Age["bob"]=18
   */
  annotations?: string;
};

export type AdminServicePullVirtualRepoParams = {
  branch?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetRepoMetaParams = {
  branch?: string;
};

export type AdminServiceUpdateServiceBody = {
  newName?: string;
};

export type AdminServiceCreateServiceParams = {
  name?: string;
};

export type AdminServiceUpdateProjectVariablesBodyVariables = {
  [key: string]: string;
};

export type AdminServiceUpdateProjectVariablesBody = {
  variables?: AdminServiceUpdateProjectVariablesBodyVariables;
};

export type AdminServiceUpdateProjectBody = {
  description?: string;
  githubUrl?: string;
  newName?: string;
  prodBranch?: string;
  prodSlots?: string;
  prodTtlSeconds?: string;
  public?: boolean;
  region?: string;
};

export type AdminServiceCreateProjectBodyVariables = { [key: string]: string };

export type AdminServiceCreateProjectBody = {
  description?: string;
  githubUrl?: string;
  name?: string;
  prodBranch?: string;
  prodOlapDriver?: string;
  prodOlapDsn?: string;
  prodSlots?: string;
  public?: boolean;
  region?: string;
  subpath?: string;
  variables?: AdminServiceCreateProjectBodyVariables;
};

export type AdminServiceListProjectsForOrganizationParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceCreateWhitelistedDomainBody = {
  domain?: string;
  role?: string;
};

export type AdminServiceSearchProjectUsersParams = {
  emailQuery?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListProjectMembersParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListProjectInvitesParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetIFrameBodyQuery = { [key: string]: string };

export type AdminServiceGetIFrameBodyAttributes = { [key: string]: any };

export type AdminServiceGetIFrameBody = {
  attributes?: AdminServiceGetIFrameBodyAttributes;
  branch?: string;
  kind?: string;
  query?: AdminServiceGetIFrameBodyQuery;
  resource?: string;
  state?: string;
  ttlSeconds?: number;
  userEmail?: string;
  userId?: string;
};

export type AdminServiceGetDeploymentCredentialsBodyAttributes = {
  [key: string]: any;
};

export type AdminServiceGetDeploymentCredentialsBody = {
  attributes?: AdminServiceGetDeploymentCredentialsBodyAttributes;
  branch?: string;
  ttlSeconds?: number;
  userEmail?: string;
  userId?: string;
};

export type AdminServiceRemoveOrganizationMemberParams = {
  keepProjectRoles?: boolean;
};

export type AdminServiceListOrganizationMembersParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListOrganizationInvitesParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceUpdateOrganizationBody = {
  description?: string;
  newName?: string;
};

export type AdminServiceListOrganizationsParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetGithubRepoStatusParams = {
  githubUrl?: string;
};

export type AdminServiceTriggerRefreshSourcesBody = {
  sources?: string[];
};

export type AdminServiceCreateReportBodyBody = {
  options?: V1ReportOptions;
};

export type AdminServiceAddOrganizationMemberBodyBody = {
  email?: string;
  role?: string;
};

export type AdminServiceSetOrganizationMemberRoleBodyBody = {
  role?: string;
};

export type AdminServiceTriggerReconcileBodyBody = { [key: string]: any };

export interface V1WhitelistedDomain {
  domain?: string;
  role?: string;
}

export interface V1VirtualFile {
  data?: string;
  deleted?: boolean;
  path?: string;
  updatedOn?: string;
}

export interface V1UserQuotas {
  singleuserOrgs?: number;
}

export interface V1UserPreferences {
  timeZone?: string;
}

export interface V1UserInvite {
  email?: string;
  invitedBy?: string;
  role?: string;
}

export interface V1User {
  createdOn?: string;
  displayName?: string;
  email?: string;
  id?: string;
  photoUrl?: string;
  quotas?: V1UserQuotas;
  updatedOn?: string;
}

export interface V1UpdateUserPreferencesResponse {
  preferences?: V1UserPreferences;
}

export interface V1UpdateUserPreferencesRequest {
  preferences?: V1UserPreferences;
}

export interface V1UpdateServiceResponse {
  service?: V1Service;
}

export type V1UpdateProjectVariablesResponseVariables = {
  [key: string]: string;
};

export interface V1UpdateProjectVariablesResponse {
  variables?: V1UpdateProjectVariablesResponseVariables;
}

export interface V1UpdateProjectResponse {
  project?: V1Project;
}

export interface V1UpdateOrganizationResponse {
  organization?: V1Organization;
}

export interface V1UnsubscribeReportResponse {
  [key: string]: any;
}

export interface V1TriggerReportResponse {
  [key: string]: any;
}

export interface V1TriggerRefreshSourcesResponse {
  [key: string]: any;
}

export interface V1TriggerRedeployResponse {
  [key: string]: any;
}

export interface V1TriggerRedeployRequest {
  deploymentId?: string;
  organization?: string;
  project?: string;
}

export interface V1TriggerReconcileResponse {
  [key: string]: any;
}

export interface V1TelemetryResponse {
  [key: string]: any;
}

export type V1TelemetryRequestEvent = { [key: string]: any };

export interface V1TelemetryRequest {
  event?: V1TelemetryRequestEvent;
  name?: string;
  value?: number;
}

export interface V1SudoUpdateUserQuotasResponse {
  user?: V1User;
}

export interface V1SudoUpdateUserQuotasRequest {
  email?: string;
  singleuserOrgs?: number;
}

export interface V1SudoUpdateOrganizationQuotasResponse {
  organization?: V1Organization;
}

export interface V1SudoUpdateOrganizationQuotasRequest {
  deployments?: number;
  orgName?: string;
  outstandingInvites?: number;
  projects?: number;
  slotsPerDeployment?: number;
  slotsTotal?: number;
}

export interface V1SudoGetResourceResponse {
  deployment?: V1Deployment;
  instance?: V1Deployment;
  org?: V1Organization;
  project?: V1Project;
  user?: V1User;
}

export interface V1SetSuperuserResponse {
  [key: string]: any;
}

export interface V1SetSuperuserRequest {
  email?: string;
  superuser?: boolean;
}

export interface V1SetProjectMemberRoleResponse {
  [key: string]: any;
}

export interface V1SetOrganizationMemberRoleResponse {
  [key: string]: any;
}

export interface V1ServiceToken {
  createdOn?: string;
  expiresOn?: string;
  id?: string;
}

export interface V1Service {
  createdOn?: string;
  id?: string;
  name?: string;
  orgId?: string;
  orgName?: string;
  updatedOn?: string;
}

export interface V1SearchUsersResponse {
  nextPageToken?: string;
  users?: V1User[];
}

export interface V1SearchProjectUsersResponse {
  nextPageToken?: string;
  users?: V1User[];
}

export interface V1SearchProjectNamesResponse {
  names?: string[];
  nextPageToken?: string;
}

export interface V1RevokeServiceAuthTokenResponse {
  [key: string]: any;
}

export interface V1RevokeCurrentAuthTokenResponse {
  tokenId?: string;
}

export interface V1ReportOptions {
  exportFormat?: V1ExportFormat;
  exportLimit?: string;
  openProjectSubpath?: string;
  queryArgsJson?: string;
  queryName?: string;
  recipients?: string[];
  refreshCron?: string;
  refreshTimeZone?: string;
  title?: string;
}

export interface V1RemoveWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1RemoveProjectMemberResponse {
  [key: string]: any;
}

export interface V1RemoveOrganizationMemberResponse {
  [key: string]: any;
}

export interface V1RemoveBookmarkResponse {
  [key: string]: any;
}

export interface V1PullVirtualRepoResponse {
  files?: V1VirtualFile[];
  nextPageToken?: string;
}

export interface V1ProjectPermissions {
  createReports?: boolean;
  manageDev?: boolean;
  manageProd?: boolean;
  manageProject?: boolean;
  manageProjectMembers?: boolean;
  manageReports?: boolean;
  readDev?: boolean;
  readDevStatus?: boolean;
  readProd?: boolean;
  readProdStatus?: boolean;
  readProject?: boolean;
  readProjectMembers?: boolean;
}

export interface V1Project {
  createdOn?: string;
  description?: string;
  frontendUrl?: string;
  githubUrl?: string;
  id?: string;
  name?: string;
  orgId?: string;
  orgName?: string;
  prodBranch?: string;
  prodDeploymentId?: string;
  prodOlapDriver?: string;
  prodOlapDsn?: string;
  prodSlots?: string;
  prodTtlSeconds?: string;
  public?: boolean;
  region?: string;
  subpath?: string;
  updatedOn?: string;
}

export interface V1PingResponse {
  time?: string;
  version?: string;
}

export interface V1OrganizationQuotas {
  deployments?: number;
  outstandingInvites?: number;
  projects?: number;
  slotsPerDeployment?: number;
  slotsTotal?: number;
}

export interface V1OrganizationPermissions {
  createProjects?: boolean;
  manageOrg?: boolean;
  manageOrgMembers?: boolean;
  manageProjects?: boolean;
  readOrg?: boolean;
  readOrgMembers?: boolean;
  readProjects?: boolean;
}

export interface V1Organization {
  createdOn?: string;
  description?: string;
  id?: string;
  name?: string;
  quotas?: V1OrganizationQuotas;
  updatedOn?: string;
}

export interface V1Member {
  createdOn?: string;
  roleName?: string;
  updatedOn?: string;
  userEmail?: string;
  userId?: string;
  userName?: string;
}

export interface V1ListWhitelistedDomainsResponse {
  domains?: V1WhitelistedDomain[];
}

export interface V1ListSuperusersResponse {
  users?: V1User[];
}

export interface V1ListServicesResponse {
  services?: V1Service[];
}

export interface V1ListServiceAuthTokensResponse {
  tokens?: V1ServiceToken[];
}

export interface V1ListProjectsForOrganizationResponse {
  nextPageToken?: string;
  projects?: V1Project[];
}

export interface V1ListProjectMembersResponse {
  members?: V1Member[];
  nextPageToken?: string;
}

export interface V1ListProjectInvitesResponse {
  invites?: V1UserInvite[];
  nextPageToken?: string;
}

export interface V1ListOrganizationsResponse {
  nextPageToken?: string;
  organizations?: V1Organization[];
}

export interface V1ListOrganizationMembersResponse {
  members?: V1Member[];
  nextPageToken?: string;
}

export interface V1ListOrganizationInvitesResponse {
  invites?: V1UserInvite[];
  nextPageToken?: string;
}

export interface V1ListBookmarksResponse {
  bookmarks?: V1Bookmark[];
}

export interface V1LeaveOrganizationResponse {
  [key: string]: any;
}

export interface V1IssueServiceAuthTokenResponse {
  token?: string;
}

export interface V1IssueRepresentativeAuthTokenResponse {
  token?: string;
}

export interface V1IssueRepresentativeAuthTokenRequest {
  email?: string;
  ttlMinutes?: string;
}

export interface V1GetUserResponse {
  user?: V1User;
}

export interface V1GetReportMetaResponse {
  editUrl?: string;
  exportUrl?: string;
  openUrl?: string;
}

export interface V1GetRepoMetaResponse {
  gitSubpath?: string;
  gitUrl?: string;
  gitUrlExpiresOn?: string;
}

export type V1GetProjectVariablesResponseVariables = { [key: string]: string };

export interface V1GetProjectVariablesResponse {
  variables?: V1GetProjectVariablesResponseVariables;
}

export interface V1GetProjectResponse {
  jwt?: string;
  prodDeployment?: V1Deployment;
  project?: V1Project;
  projectPermissions?: V1ProjectPermissions;
}

export interface V1GetOrganizationResponse {
  organization?: V1Organization;
  permissions?: V1OrganizationPermissions;
}

export interface V1GetIFrameResponse {
  accessToken?: string;
  iframeSrc?: string;
  instanceId?: string;
  runtimeHost?: string;
  ttlSeconds?: number;
}

export interface V1GetGithubRepoStatusResponse {
  defaultBranch?: string;
  grantAccessUrl?: string;
  hasAccess?: boolean;
}

export interface V1GetGitCredentialsResponse {
  password?: string;
  prodBranch?: string;
  repoUrl?: string;
  subpath?: string;
  username?: string;
}

export interface V1GetDeploymentCredentialsResponse {
  accessToken?: string;
  instanceId?: string;
  runtimeHost?: string;
  ttlSeconds?: number;
}

export interface V1GetCurrentUserResponse {
  preferences?: V1UserPreferences;
  user?: V1User;
}

export interface V1GetBookmarkResponse {
  bookmark?: V1Bookmark;
}

export interface V1GenerateReportYAMLResponse {
  yaml?: string;
}

export type V1ExportFormat =
  (typeof V1ExportFormat)[keyof typeof V1ExportFormat];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1ExportFormat = {
  EXPORT_FORMAT_UNSPECIFIED: "EXPORT_FORMAT_UNSPECIFIED",
  EXPORT_FORMAT_CSV: "EXPORT_FORMAT_CSV",
  EXPORT_FORMAT_XLSX: "EXPORT_FORMAT_XLSX",
  EXPORT_FORMAT_PARQUET: "EXPORT_FORMAT_PARQUET",
} as const;

export interface V1EditReportResponse {
  [key: string]: any;
}

export type V1DeploymentStatus =
  (typeof V1DeploymentStatus)[keyof typeof V1DeploymentStatus];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1DeploymentStatus = {
  DEPLOYMENT_STATUS_UNSPECIFIED: "DEPLOYMENT_STATUS_UNSPECIFIED",
  DEPLOYMENT_STATUS_PENDING: "DEPLOYMENT_STATUS_PENDING",
  DEPLOYMENT_STATUS_OK: "DEPLOYMENT_STATUS_OK",
  DEPLOYMENT_STATUS_ERROR: "DEPLOYMENT_STATUS_ERROR",
} as const;

export interface V1Deployment {
  branch?: string;
  createdOn?: string;
  id?: string;
  projectId?: string;
  runtimeHost?: string;
  runtimeInstanceId?: string;
  slots?: string;
  status?: V1DeploymentStatus;
  statusMessage?: string;
  updatedOn?: string;
}

export interface V1DeleteServiceResponse {
  service?: V1Service;
}

export interface V1DeleteReportResponse {
  [key: string]: any;
}

export interface V1DeleteProjectResponse {
  [key: string]: any;
}

export interface V1DeleteOrganizationResponse {
  [key: string]: any;
}

export interface V1CreateWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1CreateServiceResponse {
  service?: V1Service;
}

export interface V1CreateReportResponse {
  name?: string;
}

export interface V1CreateProjectResponse {
  project?: V1Project;
}

export interface V1CreateOrganizationResponse {
  organization?: V1Organization;
}

export interface V1CreateOrganizationRequest {
  description?: string;
  name?: string;
}

export interface V1CreateBookmarkResponse {
  bookmark?: V1Bookmark;
}

export interface V1CreateBookmarkRequest {
  dashboardName?: string;
  data?: string;
  displayName?: string;
  projectId?: string;
}

export interface V1Bookmark {
  createdOn?: string;
  dashboardName?: string;
  data?: string;
  displayName?: string;
  id?: string;
  projectId?: string;
  updatedOn?: string;
  userId?: string;
}

export interface V1AddProjectMemberResponse {
  pendingSignup?: boolean;
}

export interface V1AddOrganizationMemberResponse {
  pendingSignup?: boolean;
}

export interface RpcStatus {
  code?: number;
  details?: ProtobufAny[];
  message?: string;
}

/**
 * `NullValue` is a singleton enumeration to represent the null value for the
`Value` type union.

 The JSON representation for `NullValue` is JSON `null`.

 - NULL_VALUE: Null value.
 */
export type ProtobufNullValue =
  (typeof ProtobufNullValue)[keyof typeof ProtobufNullValue];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const ProtobufNullValue = {
  NULL_VALUE: "NULL_VALUE",
} as const;

export interface ProtobufAny {
  "@type"?: string;
  [key: string]: unknown;
}
