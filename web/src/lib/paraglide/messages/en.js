/* eslint-disable */
/** @typedef {import('../runtime.js').LocalizedString} LocalizedString */
/** @typedef {{}} Dashboard_TitleInputs */
/** @typedef {{}} App_TitleInputs */
/** @typedef {{}} Sidebar_DashboardInputs */
/** @typedef {{}} Sidebar_StatusInputs */
/** @typedef {{}} Sidebar_Status_HealthInputs */
/** @typedef {{}} Sidebar_Status_ServersInputs */
/** @typedef {{}} Sidebar_ProjectsInputs */
/** @typedef {{}} Sidebar_ServersInputs */
/** @typedef {{}} Sidebar_EventsInputs */
/** @typedef {{}} Sidebar_SignoutInputs */
/** @typedef {{}} Theme_Toggle_LightInputs */
/** @typedef {{}} Theme_Toggle_DarkInputs */
/** @typedef {{}} Login_TitleInputs */
/** @typedef {{}} Login_Email_LabelInputs */
/** @typedef {{}} Login_Password_LabelInputs */
/** @typedef {{}} Login_SubmitInputs */
/** @typedef {{}} Login_Switch_To_SignupInputs */
/** @typedef {{}} Login_Switch_To_SigninInputs */
/** @typedef {{}} Login_Signup_TitleInputs */
/** @typedef {{}} Login_Signup_SubmitInputs */
/** @typedef {{}} Login_Signup_Disabled_TitleInputs */
/** @typedef {{}} Login_Signup_Disabled_MessageInputs */
/** @typedef {{}} Login_Error_RequiredInputs */
/** @typedef {{}} Projects_TitleInputs */
/** @typedef {{}} Projects_EmptyInputs */
/** @typedef {{}} Projects_Token_StoredInputs */
/** @typedef {{}} Projects_No_TokenInputs */
/** @typedef {{}} Projects_New_Name_LabelInputs */
/** @typedef {{}} Projects_New_Token_LabelInputs */
/** @typedef {{}} Projects_New_SubmitInputs */
/** @typedef {{}} Projects_Delete_ConfirmInputs */
/** @typedef {{}} Projects_Delete_SubmitInputs */
/** @typedef {{}} Projects_Last_ErrorInputs */
/** @typedef {{}} Project_Detail_TitleInputs */
/** @typedef {{}} Project_Detail_BackInputs */
/** @typedef {{}} Project_Detail_LoadingInputs */
/** @typedef {{}} Project_Detail_Not_FoundInputs */
/** @typedef {{}} Project_Detail_Token_StoredInputs */
/** @typedef {{}} Project_Detail_Token_MissingInputs */
/** @typedef {{ date: NonNullable<unknown> }} Project_Detail_Created_AtInputs */
/** @typedef {{}} Project_Detail_Register_TitleInputs */
/** @typedef {{}} Project_Detail_Hcloud_Id_LabelInputs */
/** @typedef {{}} Project_Detail_Name_LabelInputs */
/** @typedef {{}} Project_Detail_Add_SubmitInputs */
/** @typedef {{}} Project_Detail_Add_HintInputs */
/** @typedef {{ count: NonNullable<unknown> }} Project_Detail_Servers_TitleInputs */
/** @typedef {{}} Project_Detail_Servers_EmptyInputs */
/** @typedef {{}} Servers_TitleInputs */
/** @typedef {{}} Servers_EmptyInputs */
/** @typedef {{}} Servers_Col_NameInputs */
/** @typedef {{}} Servers_Col_ProjectInputs */
/** @typedef {{}} Servers_Col_TypesInputs */
/** @typedef {{}} Servers_Col_ModeInputs */
/** @typedef {{}} Servers_Col_StatusInputs */
/** @typedef {{}} Servers_Mode_ManualInputs */
/** @typedef {{}} Servers_Mode_Auto_PromoteInputs */
/** @typedef {{}} Servers_Mode_ScheduledInputs */
/** @typedef {{}} Server_Detail_LoadingInputs */
/** @typedef {{}} Server_Detail_Not_FoundInputs */
/** @typedef {{}} Server_Detail_EditInputs */
/** @typedef {{}} Server_Detail_Tab_OverviewInputs */
/** @typedef {{}} Server_Detail_Tab_WindowsInputs */
/** @typedef {{}} Server_Detail_Tab_EventsInputs */
/** @typedef {{ id: NonNullable<unknown> }} Server_Detail_Hcloud_IdInputs */
/** @typedef {{ mode: NonNullable<unknown> }} Server_Detail_ModeInputs */
/** @typedef {{ state: NonNullable<unknown> }} Server_Detail_StateInputs */
/** @typedef {{}} Server_Detail_Base_TypeInputs */
/** @typedef {{}} Server_Detail_Top_TypeInputs */
/** @typedef {{}} Server_Detail_Fallback_ChainInputs */
/** @typedef {{}} Server_Detail_TimezoneInputs */
/** @typedef {{}} Server_Detail_Rescale_UpInputs */
/** @typedef {{}} Server_Detail_Rescale_DownInputs */
/** @typedef {{}} Server_Detail_PromoteInputs */
/** @typedef {{}} Server_Detail_DemoteInputs */
/** @typedef {{}} Server_Detail_Edit_WindowsInputs */
/** @typedef {{}} Server_Detail_Recent_EventsInputs */
/** @typedef {{}} Server_Detail_Events_EmptyInputs */
/** @typedef {{}} Server_Detail_Windows_EmptyInputs */
/** @typedef {{ count: NonNullable<unknown> }} Server_Detail_Windows_CountInputs */
/** @typedef {{}} Server_Detail_Window_EnabledInputs */
/** @typedef {{}} Server_Detail_Window_DisabledInputs */
/** @typedef {{}} Server_Detail_Window_Every_DayInputs */
/** @typedef {{}} Server_Edit_TitleInputs */
/** @typedef {{}} Server_Edit_SaveInputs */
/** @typedef {{}} Server_Edit_SavingInputs */
/** @typedef {{}} Server_Edit_CancelInputs */
/** @typedef {{}} Server_Edit_Field_NameInputs */
/** @typedef {{}} Server_Edit_Field_LabelInputs */
/** @typedef {{}} Server_Edit_Field_BaseInputs */
/** @typedef {{}} Server_Edit_Field_TopInputs */
/** @typedef {{}} Server_Edit_Field_FallbackInputs */
/** @typedef {{}} Server_Edit_Field_Fallback_PlaceholderInputs */
/** @typedef {{}} Server_Edit_Field_ModeInputs */
/** @typedef {{}} Server_Edit_Field_TimezoneInputs */
/** @typedef {{}} Server_Edit_Field_Timezone_PlaceholderInputs */
/** @typedef {{}} Windows_TitleInputs */
/** @typedef {{}} Windows_AddInputs */
/** @typedef {{}} Windows_EmptyInputs */
/** @typedef {{}} Windows_Col_LabelInputs */
/** @typedef {{}} Windows_Col_DaysInputs */
/** @typedef {{}} Windows_Col_StartInputs */
/** @typedef {{}} Windows_Col_StopInputs */
/** @typedef {{}} Windows_Col_TargetInputs */
/** @typedef {{}} Windows_Col_EnabledInputs */
/** @typedef {{}} Windows_Col_YesInputs */
/** @typedef {{}} Windows_Col_NoInputs */
/** @typedef {{}} Windows_EnableInputs */
/** @typedef {{}} Windows_DisableInputs */
/** @typedef {{}} Windows_DeleteInputs */
/** @typedef {{}} Windows_Delete_ConfirmInputs */
/** @typedef {{}} Windows_Modal_TitleInputs */
/** @typedef {{}} Windows_Modal_SaveInputs */
/** @typedef {{}} Windows_Modal_SavingInputs */
/** @typedef {{}} Windows_Modal_CancelInputs */
/** @typedef {{}} Windows_Field_LabelInputs */
/** @typedef {{}} Windows_Field_DaysInputs */
/** @typedef {{}} Windows_Field_StartInputs */
/** @typedef {{}} Windows_Field_StopInputs */
/** @typedef {{}} Windows_Field_TargetInputs */
/** @typedef {{}} Windows_Field_EnabledInputs */
/** @typedef {{}} Events_TitleInputs */
/** @typedef {{}} Events_Filter_ServerInputs */
/** @typedef {{}} Events_Filter_Server_AllInputs */
/** @typedef {{}} Events_Filter_KindInputs */
/** @typedef {{}} Events_Filter_Kind_AllInputs */
/** @typedef {{}} Events_Filter_LimitInputs */
/** @typedef {{}} Events_LoadingInputs */
/** @typedef {{}} Events_EmptyInputs */
/** @typedef {{}} Dashboard_LoadingInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Section_ProjectsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Section_ServersInputs */
/** @typedef {{}} Dashboard_Section_Recent_EventsInputs */
/** @typedef {{}} Kpi_Active_ServersInputs */
/** @typedef {{}} Kpi_Active_Servers_HintInputs */
/** @typedef {{}} Kpi_ProjectsInputs */
/** @typedef {{}} Kpi_Projects_HintInputs */
/** @typedef {{}} Kpi_Rescales_24h_OkInputs */
/** @typedef {{}} Kpi_Last_ErrorInputs */
/** @typedef {{}} Kpi_No_ErrorInputs */
/** @typedef {{}} Kpi_LoadingInputs */
/** @typedef {{}} Dashboard_Chart_ActivityInputs */
/** @typedef {{}} Dashboard_Chart_CostInputs */
/** @typedef {{}} Dashboard_Chart_RangeInputs */
/** @typedef {{}} Dashboard_Chart_Range_1dInputs */
/** @typedef {{}} Dashboard_Chart_Range_7dInputs */
/** @typedef {{}} Dashboard_Chart_Range_30dInputs */
/** @typedef {{}} Dashboard_Chart_Cost_EmptyInputs */
/** @typedef {{}} Health_TitleInputs */
/** @typedef {{}} Health_Card_Api_LabelInputs */
/** @typedef {{}} Health_Card_Db_LabelInputs */
/** @typedef {{}} Health_Card_Hcloud_LabelInputs */
/** @typedef {{}} Health_Card_Last_Event_LabelInputs */
/** @typedef {{}} Health_Card_Recent_Errors_LabelInputs */
/** @typedef {{}} Health_Card_Windows_LabelInputs */
/** @typedef {{ threshold: NonNullable<unknown> }} Health_Warn_ThresholdInputs */
/** @typedef {{}} Health_Ok_BelowInputs */
/** @typedef {{}} Health_Warn_AboveInputs */
/** @typedef {{}} Health_Fail_AboveInputs */
/** @typedef {{}} Health_CheckingInputs */
/** @typedef {{}} Servers_Status_TitleInputs */
/** @typedef {{}} Servers_Status_EmptyInputs */
/** @typedef {{}} Servers_Status_Col_NameInputs */
/** @typedef {{}} Servers_Status_Col_ModeInputs */
/** @typedef {{}} Servers_Status_Col_TopInputs */
/** @typedef {{}} Servers_Status_Col_WindowInputs */
/** @typedef {{}} Servers_Status_Col_StatusInputs */


export const dashboard_title = /** @type {(inputs: Dashboard_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dashboard`)
};

export const app_title = /** @type {(inputs: App_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Hetzner Rescaler`)
};

export const sidebar_dashboard = /** @type {(inputs: Sidebar_DashboardInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dashboard`)
};

export const sidebar_status = /** @type {(inputs: Sidebar_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Status`)
};

export const sidebar_status_health = /** @type {(inputs: Sidebar_Status_HealthInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Health`)
};

export const sidebar_status_servers = /** @type {(inputs: Sidebar_Status_ServersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Operational`)
};

export const sidebar_projects = /** @type {(inputs: Sidebar_ProjectsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Projects`)
};

export const sidebar_servers = /** @type {(inputs: Sidebar_ServersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Servers (registered)`)
};

export const sidebar_events = /** @type {(inputs: Sidebar_EventsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Events`)
};

export const sidebar_signout = /** @type {(inputs: Sidebar_SignoutInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sign out`)
};

export const theme_toggle_light = /** @type {(inputs: Theme_Toggle_LightInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Switch to light mode`)
};

export const theme_toggle_dark = /** @type {(inputs: Theme_Toggle_DarkInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Switch to dark mode`)
};

export const login_title = /** @type {(inputs: Login_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sign in`)
};

export const login_email_label = /** @type {(inputs: Login_Email_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Email`)
};

export const login_password_label = /** @type {(inputs: Login_Password_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Password`)
};

export const login_submit = /** @type {(inputs: Login_SubmitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sign in`)
};

export const login_switch_to_signup = /** @type {(inputs: Login_Switch_To_SignupInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Need an account? Sign up`)
};

export const login_switch_to_signin = /** @type {(inputs: Login_Switch_To_SigninInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Already have an account? Sign in`)
};

export const login_signup_title = /** @type {(inputs: Login_Signup_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Create account`)
};

export const login_signup_submit = /** @type {(inputs: Login_Signup_SubmitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Create account`)
};

export const login_signup_disabled_title = /** @type {(inputs: Login_Signup_Disabled_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Signups are closed`)
};

export const login_signup_disabled_message = /** @type {(inputs: Login_Signup_Disabled_MessageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`New account registration has been disabled by the operator. If you need an account, ask an existing user to add one.`)
};

export const login_error_required = /** @type {(inputs: Login_Error_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Email and password are required`)
};

export const projects_title = /** @type {(inputs: Projects_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Projects`)
};

export const projects_empty = /** @type {(inputs: Projects_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No projects yet. Create one below.`)
};

export const projects_token_stored = /** @type {(inputs: Projects_Token_StoredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`token stored`)
};

export const projects_no_token = /** @type {(inputs: Projects_No_TokenInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`no token`)
};

export const projects_new_name_label = /** @type {(inputs: Projects_New_Name_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Project name`)
};

export const projects_new_token_label = /** @type {(inputs: Projects_New_Token_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Hetzner API token`)
};

export const projects_new_submit = /** @type {(inputs: Projects_New_SubmitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Create project`)
};

export const projects_delete_confirm = /** @type {(inputs: Projects_Delete_ConfirmInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Delete project and all its servers?`)
};

export const projects_delete_submit = /** @type {(inputs: Projects_Delete_SubmitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Delete`)
};

export const projects_last_error = /** @type {(inputs: Projects_Last_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Last error`)
};

export const project_detail_title = /** @type {(inputs: Project_Detail_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Project`)
};

export const project_detail_back = /** @type {(inputs: Project_Detail_BackInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Back to projects`)
};

export const project_detail_loading = /** @type {(inputs: Project_Detail_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Loading…`)
};

export const project_detail_not_found = /** @type {(inputs: Project_Detail_Not_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Project not found.`)
};

export const project_detail_token_stored = /** @type {(inputs: Project_Detail_Token_StoredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`token stored`)
};

export const project_detail_token_missing = /** @type {(inputs: Project_Detail_Token_MissingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`missing`)
};

export const project_detail_created_at = /** @type {(inputs: Project_Detail_Created_AtInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`created ${i?.date}`)
};

export const project_detail_register_title = /** @type {(inputs: Project_Detail_Register_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Register a server manually`)
};

export const project_detail_hcloud_id_label = /** @type {(inputs: Project_Detail_Hcloud_Id_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Hetzner server ID`)
};

export const project_detail_name_label = /** @type {(inputs: Project_Detail_Name_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Display name`)
};

export const project_detail_add_submit = /** @type {(inputs: Project_Detail_Add_SubmitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Add server`)
};

export const project_detail_add_hint = /** @type {(inputs: Project_Detail_Add_HintInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Default base/top types are filled in. Edit them on the server detail page.`)
};

export const project_detail_servers_title = /** @type {(inputs: Project_Detail_Servers_TitleInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Servers (${i?.count})`)
};

export const project_detail_servers_empty = /** @type {(inputs: Project_Detail_Servers_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No servers registered.`)
};

export const servers_title = /** @type {(inputs: Servers_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Servers`)
};

export const servers_empty = /** @type {(inputs: Servers_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No servers registered.`)
};

export const servers_col_name = /** @type {(inputs: Servers_Col_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Name`)
};

export const servers_col_project = /** @type {(inputs: Servers_Col_ProjectInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Project`)
};

export const servers_col_types = /** @type {(inputs: Servers_Col_TypesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Base → Top`)
};

export const servers_col_mode = /** @type {(inputs: Servers_Col_ModeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mode`)
};

export const servers_col_status = /** @type {(inputs: Servers_Col_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Status`)
};

export const servers_mode_manual = /** @type {(inputs: Servers_Mode_ManualInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Manual`)
};

export const servers_mode_auto_promote = /** @type {(inputs: Servers_Mode_Auto_PromoteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Auto-promote`)
};

export const servers_mode_scheduled = /** @type {(inputs: Servers_Mode_ScheduledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Scheduled`)
};

export const server_detail_loading = /** @type {(inputs: Server_Detail_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Loading…`)
};

export const server_detail_not_found = /** @type {(inputs: Server_Detail_Not_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Server not found.`)
};

export const server_detail_edit = /** @type {(inputs: Server_Detail_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Edit`)
};

export const server_detail_tab_overview = /** @type {(inputs: Server_Detail_Tab_OverviewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Overview`)
};

export const server_detail_tab_windows = /** @type {(inputs: Server_Detail_Tab_WindowsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Windows`)
};

export const server_detail_tab_events = /** @type {(inputs: Server_Detail_Tab_EventsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Events`)
};

export const server_detail_hcloud_id = /** @type {(inputs: Server_Detail_Hcloud_IdInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Hetzner #${i?.id}`)
};

export const server_detail_mode = /** @type {(inputs: Server_Detail_ModeInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`mode: ${i?.mode}`)
};

export const server_detail_state = /** @type {(inputs: Server_Detail_StateInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`state: ${i?.state}`)
};

export const server_detail_base_type = /** @type {(inputs: Server_Detail_Base_TypeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Base type`)
};

export const server_detail_top_type = /** @type {(inputs: Server_Detail_Top_TypeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Top type`)
};

export const server_detail_fallback_chain = /** @type {(inputs: Server_Detail_Fallback_ChainInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fallback chain`)
};

export const server_detail_timezone = /** @type {(inputs: Server_Detail_TimezoneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Timezone`)
};

export const server_detail_rescale_up = /** @type {(inputs: Server_Detail_Rescale_UpInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rescale up`)
};

export const server_detail_rescale_down = /** @type {(inputs: Server_Detail_Rescale_DownInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rescale down`)
};

export const server_detail_promote = /** @type {(inputs: Server_Detail_PromoteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Request promote`)
};

export const server_detail_demote = /** @type {(inputs: Server_Detail_DemoteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Request demote`)
};

export const server_detail_edit_windows = /** @type {(inputs: Server_Detail_Edit_WindowsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Edit windows`)
};

export const server_detail_recent_events = /** @type {(inputs: Server_Detail_Recent_EventsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Recent events`)
};

export const server_detail_events_empty = /** @type {(inputs: Server_Detail_Events_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No events yet.`)
};

export const server_detail_windows_empty = /** @type {(inputs: Server_Detail_Windows_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No windows yet.`)
};

export const server_detail_windows_count = /** @type {(inputs: Server_Detail_Windows_CountInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Existing windows (${i?.count})`)
};

export const server_detail_window_enabled = /** @type {(inputs: Server_Detail_Window_EnabledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Enabled`)
};

export const server_detail_window_disabled = /** @type {(inputs: Server_Detail_Window_DisabledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Disabled`)
};

export const server_detail_window_every_day = /** @type {(inputs: Server_Detail_Window_Every_DayInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Every day`)
};

export const server_edit_title = /** @type {(inputs: Server_Edit_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Edit server`)
};

export const server_edit_save = /** @type {(inputs: Server_Edit_SaveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Save`)
};

export const server_edit_saving = /** @type {(inputs: Server_Edit_SavingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Saving…`)
};

export const server_edit_cancel = /** @type {(inputs: Server_Edit_CancelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cancel`)
};

export const server_edit_field_name = /** @type {(inputs: Server_Edit_Field_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Name`)
};

export const server_edit_field_label = /** @type {(inputs: Server_Edit_Field_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Label`)
};

export const server_edit_field_base = /** @type {(inputs: Server_Edit_Field_BaseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Base server type`)
};

export const server_edit_field_top = /** @type {(inputs: Server_Edit_Field_TopInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Top server type`)
};

export const server_edit_field_fallback = /** @type {(inputs: Server_Edit_Field_FallbackInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fallback chain (comma-separated, top first)`)
};

export const server_edit_field_fallback_placeholder = /** @type {(inputs: Server_Edit_Field_Fallback_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`cpx31, cpx21, cpx11`)
};

export const server_edit_field_mode = /** @type {(inputs: Server_Edit_Field_ModeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mode`)
};

export const server_edit_field_timezone = /** @type {(inputs: Server_Edit_Field_TimezoneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Timezone (IANA)`)
};

export const server_edit_field_timezone_placeholder = /** @type {(inputs: Server_Edit_Field_Timezone_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Europe/Rome`)
};

export const windows_title = /** @type {(inputs: Windows_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Windows`)
};

export const windows_add = /** @type {(inputs: Windows_AddInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Add window`)
};

export const windows_empty = /** @type {(inputs: Windows_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No windows yet.`)
};

export const windows_col_label = /** @type {(inputs: Windows_Col_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Label`)
};

export const windows_col_days = /** @type {(inputs: Windows_Col_DaysInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Days`)
};

export const windows_col_start = /** @type {(inputs: Windows_Col_StartInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Start`)
};

export const windows_col_stop = /** @type {(inputs: Windows_Col_StopInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stop`)
};

export const windows_col_target = /** @type {(inputs: Windows_Col_TargetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Target type`)
};

export const windows_col_enabled = /** @type {(inputs: Windows_Col_EnabledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Enabled`)
};

export const windows_col_yes = /** @type {(inputs: Windows_Col_YesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Yes`)
};

export const windows_col_no = /** @type {(inputs: Windows_Col_NoInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No`)
};

export const windows_enable = /** @type {(inputs: Windows_EnableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Enable`)
};

export const windows_disable = /** @type {(inputs: Windows_DisableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Disable`)
};

export const windows_delete = /** @type {(inputs: Windows_DeleteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Delete`)
};

export const windows_delete_confirm = /** @type {(inputs: Windows_Delete_ConfirmInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Delete this window?`)
};

export const windows_modal_title = /** @type {(inputs: Windows_Modal_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`New window`)
};

export const windows_modal_save = /** @type {(inputs: Windows_Modal_SaveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Create`)
};

export const windows_modal_saving = /** @type {(inputs: Windows_Modal_SavingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creating…`)
};

export const windows_modal_cancel = /** @type {(inputs: Windows_Modal_CancelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cancel`)
};

export const windows_field_label = /** @type {(inputs: Windows_Field_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Label`)
};

export const windows_field_days = /** @type {(inputs: Windows_Field_DaysInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Days of week`)
};

export const windows_field_start = /** @type {(inputs: Windows_Field_StartInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Start time`)
};

export const windows_field_stop = /** @type {(inputs: Windows_Field_StopInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stop time`)
};

export const windows_field_target = /** @type {(inputs: Windows_Field_TargetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Target server type`)
};

export const windows_field_enabled = /** @type {(inputs: Windows_Field_EnabledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Enabled`)
};

export const events_title = /** @type {(inputs: Events_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Events`)
};

export const events_filter_server = /** @type {(inputs: Events_Filter_ServerInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Server`)
};

export const events_filter_server_all = /** @type {(inputs: Events_Filter_Server_AllInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`All servers`)
};

export const events_filter_kind = /** @type {(inputs: Events_Filter_KindInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Kind`)
};

export const events_filter_kind_all = /** @type {(inputs: Events_Filter_Kind_AllInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`All kinds`)
};

export const events_filter_limit = /** @type {(inputs: Events_Filter_LimitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Limit`)
};

export const events_loading = /** @type {(inputs: Events_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Loading…`)
};

export const events_empty = /** @type {(inputs: Events_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No events yet.`)
};

export const dashboard_loading = /** @type {(inputs: Dashboard_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Loading…`)
};

export const dashboard_section_projects = /** @type {(inputs: Dashboard_Section_ProjectsInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Projects (${i?.count})`)
};

export const dashboard_section_servers = /** @type {(inputs: Dashboard_Section_ServersInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Servers (${i?.count})`)
};

export const dashboard_section_recent_events = /** @type {(inputs: Dashboard_Section_Recent_EventsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Recent events`)
};

export const kpi_active_servers = /** @type {(inputs: Kpi_Active_ServersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Active servers`)
};

export const kpi_active_servers_hint = /** @type {(inputs: Kpi_Active_Servers_HintInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Servers in auto_promote or scheduled mode`)
};

export const kpi_projects = /** @type {(inputs: Kpi_ProjectsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Projects`)
};

export const kpi_projects_hint = /** @type {(inputs: Kpi_Projects_HintInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Projects with a stored Hetzner API token`)
};

export const kpi_rescales_24h_ok = /** @type {(inputs: Kpi_Rescales_24h_OkInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rescales (24h, successful)`)
};

export const kpi_last_error = /** @type {(inputs: Kpi_Last_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Last rescale error`)
};

export const kpi_no_error = /** @type {(inputs: Kpi_No_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No recent failures`)
};

export const kpi_loading = /** @type {(inputs: Kpi_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Loading…`)
};

export const dashboard_chart_activity = /** @type {(inputs: Dashboard_Chart_ActivityInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rescaling activity`)
};

export const dashboard_chart_cost = /** @type {(inputs: Dashboard_Chart_CostInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cost breakdown (€ / month)`)
};

export const dashboard_chart_range = /** @type {(inputs: Dashboard_Chart_RangeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Range`)
};

export const dashboard_chart_range_1d = /** @type {(inputs: Dashboard_Chart_Range_1dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`1d`)
};

export const dashboard_chart_range_7d = /** @type {(inputs: Dashboard_Chart_Range_7dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`7d`)
};

export const dashboard_chart_range_30d = /** @type {(inputs: Dashboard_Chart_Range_30dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`30d`)
};

export const dashboard_chart_cost_empty = /** @type {(inputs: Dashboard_Chart_Cost_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No rescaling history in this range.`)
};

export const health_title = /** @type {(inputs: Health_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`System health`)
};

export const health_card_api_label = /** @type {(inputs: Health_Card_Api_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Go API reachable`)
};

export const health_card_db_label = /** @type {(inputs: Health_Card_Db_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`SQLite writable`)
};

export const health_card_hcloud_label = /** @type {(inputs: Health_Card_Hcloud_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Hetzner API reachable`)
};

export const health_card_last_event_label = /** @type {(inputs: Health_Card_Last_Event_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Last event age`)
};

export const health_card_recent_errors_label = /** @type {(inputs: Health_Card_Recent_Errors_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Recent errors (24h)`)
};

export const health_card_windows_label = /** @type {(inputs: Health_Card_Windows_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Active rescale windows`)
};

export const health_warn_threshold = /** @type {(inputs: Health_Warn_ThresholdInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Warning threshold: ${i?.threshold}`)
};

export const health_ok_below = /** @type {(inputs: Health_Ok_BelowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Healthy`)
};

export const health_warn_above = /** @type {(inputs: Health_Warn_AboveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Warning`)
};

export const health_fail_above = /** @type {(inputs: Health_Fail_AboveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Failing`)
};

export const health_checking = /** @type {(inputs: Health_CheckingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Refresh`)
};

export const servers_status_title = /** @type {(inputs: Servers_Status_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Server status`)
};

export const servers_status_empty = /** @type {(inputs: Servers_Status_EmptyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`No servers registered.`)
};

export const servers_status_col_name = /** @type {(inputs: Servers_Status_Col_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Server`)
};

export const servers_status_col_mode = /** @type {(inputs: Servers_Status_Col_ModeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mode`)
};

export const servers_status_col_top = /** @type {(inputs: Servers_Status_Col_TopInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Top type`)
};

export const servers_status_col_window = /** @type {(inputs: Servers_Status_Col_WindowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Active window`)
};

export const servers_status_col_status = /** @type {(inputs: Servers_Status_Col_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Status`)
};