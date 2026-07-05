/* eslint-disable */
import { getLocale, experimentalStaticLocale } from "../runtime.js"

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
import * as __en from "./en.js"
/**
* | output |
* | --- |
* | "Dashboard" |
*
* @param {Dashboard_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_title = /** @type {((inputs?: Dashboard_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_title(inputs)
});
/**
* | output |
* | --- |
* | "Hetzner Rescaler" |
*
* @param {App_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const app_title = /** @type {((inputs?: App_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<App_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.app_title(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard" |
*
* @param {Sidebar_DashboardInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_dashboard = /** @type {((inputs?: Sidebar_DashboardInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_DashboardInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Sidebar_StatusInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_status = /** @type {((inputs?: Sidebar_StatusInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_StatusInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_status(inputs)
});
/**
* | output |
* | --- |
* | "Health" |
*
* @param {Sidebar_Status_HealthInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_status_health = /** @type {((inputs?: Sidebar_Status_HealthInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_Status_HealthInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_status_health(inputs)
});
/**
* | output |
* | --- |
* | "Operational" |
*
* @param {Sidebar_Status_ServersInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_status_servers = /** @type {((inputs?: Sidebar_Status_ServersInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_Status_ServersInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_status_servers(inputs)
});
/**
* | output |
* | --- |
* | "Projects" |
*
* @param {Sidebar_ProjectsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_projects = /** @type {((inputs?: Sidebar_ProjectsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_ProjectsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_projects(inputs)
});
/**
* | output |
* | --- |
* | "Servers (registered)" |
*
* @param {Sidebar_ServersInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_servers = /** @type {((inputs?: Sidebar_ServersInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_ServersInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_servers(inputs)
});
/**
* | output |
* | --- |
* | "Events" |
*
* @param {Sidebar_EventsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_events = /** @type {((inputs?: Sidebar_EventsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_EventsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_events(inputs)
});
/**
* | output |
* | --- |
* | "Sign out" |
*
* @param {Sidebar_SignoutInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const sidebar_signout = /** @type {((inputs?: Sidebar_SignoutInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_SignoutInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.sidebar_signout(inputs)
});
/**
* | output |
* | --- |
* | "Switch to light mode" |
*
* @param {Theme_Toggle_LightInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const theme_toggle_light = /** @type {((inputs?: Theme_Toggle_LightInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_Toggle_LightInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.theme_toggle_light(inputs)
});
/**
* | output |
* | --- |
* | "Switch to dark mode" |
*
* @param {Theme_Toggle_DarkInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const theme_toggle_dark = /** @type {((inputs?: Theme_Toggle_DarkInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_Toggle_DarkInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.theme_toggle_dark(inputs)
});
/**
* | output |
* | --- |
* | "Sign in" |
*
* @param {Login_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_title = /** @type {((inputs?: Login_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_title(inputs)
});
/**
* | output |
* | --- |
* | "Email" |
*
* @param {Login_Email_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_email_label = /** @type {((inputs?: Login_Email_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Email_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_email_label(inputs)
});
/**
* | output |
* | --- |
* | "Password" |
*
* @param {Login_Password_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_password_label = /** @type {((inputs?: Login_Password_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Password_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_password_label(inputs)
});
/**
* | output |
* | --- |
* | "Sign in" |
*
* @param {Login_SubmitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_submit = /** @type {((inputs?: Login_SubmitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_SubmitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_submit(inputs)
});
/**
* | output |
* | --- |
* | "Need an account? Sign up" |
*
* @param {Login_Switch_To_SignupInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_switch_to_signup = /** @type {((inputs?: Login_Switch_To_SignupInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Switch_To_SignupInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_switch_to_signup(inputs)
});
/**
* | output |
* | --- |
* | "Already have an account? Sign in" |
*
* @param {Login_Switch_To_SigninInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_switch_to_signin = /** @type {((inputs?: Login_Switch_To_SigninInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Switch_To_SigninInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_switch_to_signin(inputs)
});
/**
* | output |
* | --- |
* | "Create account" |
*
* @param {Login_Signup_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_signup_title = /** @type {((inputs?: Login_Signup_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Signup_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_signup_title(inputs)
});
/**
* | output |
* | --- |
* | "Create account" |
*
* @param {Login_Signup_SubmitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_signup_submit = /** @type {((inputs?: Login_Signup_SubmitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Signup_SubmitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_signup_submit(inputs)
});
/**
* | output |
* | --- |
* | "Signups are closed" |
*
* @param {Login_Signup_Disabled_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_signup_disabled_title = /** @type {((inputs?: Login_Signup_Disabled_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Signup_Disabled_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_signup_disabled_title(inputs)
});
/**
* | output |
* | --- |
* | "New account registration has been disabled by the operator. If you need an account, ask an existing user to add one." |
*
* @param {Login_Signup_Disabled_MessageInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_signup_disabled_message = /** @type {((inputs?: Login_Signup_Disabled_MessageInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Signup_Disabled_MessageInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_signup_disabled_message(inputs)
});
/**
* | output |
* | --- |
* | "Email and password are required" |
*
* @param {Login_Error_RequiredInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const login_error_required = /** @type {((inputs?: Login_Error_RequiredInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Login_Error_RequiredInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.login_error_required(inputs)
});
/**
* | output |
* | --- |
* | "Projects" |
*
* @param {Projects_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_title = /** @type {((inputs?: Projects_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_title(inputs)
});
/**
* | output |
* | --- |
* | "No projects yet. Create one below." |
*
* @param {Projects_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_empty = /** @type {((inputs?: Projects_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_empty(inputs)
});
/**
* | output |
* | --- |
* | "token stored" |
*
* @param {Projects_Token_StoredInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_token_stored = /** @type {((inputs?: Projects_Token_StoredInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_Token_StoredInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_token_stored(inputs)
});
/**
* | output |
* | --- |
* | "no token" |
*
* @param {Projects_No_TokenInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_no_token = /** @type {((inputs?: Projects_No_TokenInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_No_TokenInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_no_token(inputs)
});
/**
* | output |
* | --- |
* | "Project name" |
*
* @param {Projects_New_Name_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_new_name_label = /** @type {((inputs?: Projects_New_Name_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_New_Name_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_new_name_label(inputs)
});
/**
* | output |
* | --- |
* | "Hetzner API token" |
*
* @param {Projects_New_Token_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_new_token_label = /** @type {((inputs?: Projects_New_Token_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_New_Token_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_new_token_label(inputs)
});
/**
* | output |
* | --- |
* | "Create project" |
*
* @param {Projects_New_SubmitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_new_submit = /** @type {((inputs?: Projects_New_SubmitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_New_SubmitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_new_submit(inputs)
});
/**
* | output |
* | --- |
* | "Delete project and all its servers?" |
*
* @param {Projects_Delete_ConfirmInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_delete_confirm = /** @type {((inputs?: Projects_Delete_ConfirmInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_Delete_ConfirmInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_delete_confirm(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Projects_Delete_SubmitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_delete_submit = /** @type {((inputs?: Projects_Delete_SubmitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_Delete_SubmitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_delete_submit(inputs)
});
/**
* | output |
* | --- |
* | "Last error" |
*
* @param {Projects_Last_ErrorInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const projects_last_error = /** @type {((inputs?: Projects_Last_ErrorInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Projects_Last_ErrorInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.projects_last_error(inputs)
});
/**
* | output |
* | --- |
* | "Project" |
*
* @param {Project_Detail_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_title = /** @type {((inputs?: Project_Detail_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_title(inputs)
});
/**
* | output |
* | --- |
* | "Back to projects" |
*
* @param {Project_Detail_BackInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_back = /** @type {((inputs?: Project_Detail_BackInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_BackInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_back(inputs)
});
/**
* | output |
* | --- |
* | "Loading…" |
*
* @param {Project_Detail_LoadingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_loading = /** @type {((inputs?: Project_Detail_LoadingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_LoadingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_loading(inputs)
});
/**
* | output |
* | --- |
* | "Project not found." |
*
* @param {Project_Detail_Not_FoundInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_not_found = /** @type {((inputs?: Project_Detail_Not_FoundInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Not_FoundInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_not_found(inputs)
});
/**
* | output |
* | --- |
* | "token stored" |
*
* @param {Project_Detail_Token_StoredInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_token_stored = /** @type {((inputs?: Project_Detail_Token_StoredInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Token_StoredInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_token_stored(inputs)
});
/**
* | output |
* | --- |
* | "missing" |
*
* @param {Project_Detail_Token_MissingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_token_missing = /** @type {((inputs?: Project_Detail_Token_MissingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Token_MissingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_token_missing(inputs)
});
/**
* | output |
* | --- |
* | "created {date}" |
*
* @param {Project_Detail_Created_AtInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_created_at = /** @type {((inputs: Project_Detail_Created_AtInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Created_AtInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_created_at(inputs)
});
/**
* | output |
* | --- |
* | "Register a server manually" |
*
* @param {Project_Detail_Register_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_register_title = /** @type {((inputs?: Project_Detail_Register_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Register_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_register_title(inputs)
});
/**
* | output |
* | --- |
* | "Hetzner server ID" |
*
* @param {Project_Detail_Hcloud_Id_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_hcloud_id_label = /** @type {((inputs?: Project_Detail_Hcloud_Id_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Hcloud_Id_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_hcloud_id_label(inputs)
});
/**
* | output |
* | --- |
* | "Display name" |
*
* @param {Project_Detail_Name_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_name_label = /** @type {((inputs?: Project_Detail_Name_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Name_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_name_label(inputs)
});
/**
* | output |
* | --- |
* | "Add server" |
*
* @param {Project_Detail_Add_SubmitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_add_submit = /** @type {((inputs?: Project_Detail_Add_SubmitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Add_SubmitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_add_submit(inputs)
});
/**
* | output |
* | --- |
* | "Default base/top types are filled in. Edit them on the server detail page." |
*
* @param {Project_Detail_Add_HintInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_add_hint = /** @type {((inputs?: Project_Detail_Add_HintInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Add_HintInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_add_hint(inputs)
});
/**
* | output |
* | --- |
* | "Servers ({count})" |
*
* @param {Project_Detail_Servers_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_servers_title = /** @type {((inputs: Project_Detail_Servers_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Servers_TitleInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_servers_title(inputs)
});
/**
* | output |
* | --- |
* | "No servers registered." |
*
* @param {Project_Detail_Servers_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const project_detail_servers_empty = /** @type {((inputs?: Project_Detail_Servers_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Detail_Servers_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.project_detail_servers_empty(inputs)
});
/**
* | output |
* | --- |
* | "Servers" |
*
* @param {Servers_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_title = /** @type {((inputs?: Servers_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_title(inputs)
});
/**
* | output |
* | --- |
* | "No servers registered." |
*
* @param {Servers_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_empty = /** @type {((inputs?: Servers_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_empty(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Servers_Col_NameInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_col_name = /** @type {((inputs?: Servers_Col_NameInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Col_NameInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_col_name(inputs)
});
/**
* | output |
* | --- |
* | "Project" |
*
* @param {Servers_Col_ProjectInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_col_project = /** @type {((inputs?: Servers_Col_ProjectInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Col_ProjectInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_col_project(inputs)
});
/**
* | output |
* | --- |
* | "Base → Top" |
*
* @param {Servers_Col_TypesInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_col_types = /** @type {((inputs?: Servers_Col_TypesInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Col_TypesInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_col_types(inputs)
});
/**
* | output |
* | --- |
* | "Mode" |
*
* @param {Servers_Col_ModeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_col_mode = /** @type {((inputs?: Servers_Col_ModeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Col_ModeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_col_mode(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Servers_Col_StatusInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_col_status = /** @type {((inputs?: Servers_Col_StatusInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Col_StatusInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_col_status(inputs)
});
/**
* | output |
* | --- |
* | "Manual" |
*
* @param {Servers_Mode_ManualInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_mode_manual = /** @type {((inputs?: Servers_Mode_ManualInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Mode_ManualInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_mode_manual(inputs)
});
/**
* | output |
* | --- |
* | "Auto-promote" |
*
* @param {Servers_Mode_Auto_PromoteInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_mode_auto_promote = /** @type {((inputs?: Servers_Mode_Auto_PromoteInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Mode_Auto_PromoteInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_mode_auto_promote(inputs)
});
/**
* | output |
* | --- |
* | "Scheduled" |
*
* @param {Servers_Mode_ScheduledInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_mode_scheduled = /** @type {((inputs?: Servers_Mode_ScheduledInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Mode_ScheduledInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_mode_scheduled(inputs)
});
/**
* | output |
* | --- |
* | "Loading…" |
*
* @param {Server_Detail_LoadingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_loading = /** @type {((inputs?: Server_Detail_LoadingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_LoadingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_loading(inputs)
});
/**
* | output |
* | --- |
* | "Server not found." |
*
* @param {Server_Detail_Not_FoundInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_not_found = /** @type {((inputs?: Server_Detail_Not_FoundInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Not_FoundInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_not_found(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Server_Detail_EditInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_edit = /** @type {((inputs?: Server_Detail_EditInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_EditInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_edit(inputs)
});
/**
* | output |
* | --- |
* | "Overview" |
*
* @param {Server_Detail_Tab_OverviewInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_tab_overview = /** @type {((inputs?: Server_Detail_Tab_OverviewInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Tab_OverviewInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_tab_overview(inputs)
});
/**
* | output |
* | --- |
* | "Windows" |
*
* @param {Server_Detail_Tab_WindowsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_tab_windows = /** @type {((inputs?: Server_Detail_Tab_WindowsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Tab_WindowsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_tab_windows(inputs)
});
/**
* | output |
* | --- |
* | "Events" |
*
* @param {Server_Detail_Tab_EventsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_tab_events = /** @type {((inputs?: Server_Detail_Tab_EventsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Tab_EventsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_tab_events(inputs)
});
/**
* | output |
* | --- |
* | "Hetzner #{id}" |
*
* @param {Server_Detail_Hcloud_IdInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_hcloud_id = /** @type {((inputs: Server_Detail_Hcloud_IdInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Hcloud_IdInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_hcloud_id(inputs)
});
/**
* | output |
* | --- |
* | "mode: {mode}" |
*
* @param {Server_Detail_ModeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_mode = /** @type {((inputs: Server_Detail_ModeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_ModeInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_mode(inputs)
});
/**
* | output |
* | --- |
* | "state: {state}" |
*
* @param {Server_Detail_StateInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_state = /** @type {((inputs: Server_Detail_StateInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_StateInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_state(inputs)
});
/**
* | output |
* | --- |
* | "Base type" |
*
* @param {Server_Detail_Base_TypeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_base_type = /** @type {((inputs?: Server_Detail_Base_TypeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Base_TypeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_base_type(inputs)
});
/**
* | output |
* | --- |
* | "Top type" |
*
* @param {Server_Detail_Top_TypeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_top_type = /** @type {((inputs?: Server_Detail_Top_TypeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Top_TypeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_top_type(inputs)
});
/**
* | output |
* | --- |
* | "Fallback chain" |
*
* @param {Server_Detail_Fallback_ChainInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_fallback_chain = /** @type {((inputs?: Server_Detail_Fallback_ChainInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Fallback_ChainInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_fallback_chain(inputs)
});
/**
* | output |
* | --- |
* | "Timezone" |
*
* @param {Server_Detail_TimezoneInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_timezone = /** @type {((inputs?: Server_Detail_TimezoneInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_TimezoneInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_timezone(inputs)
});
/**
* | output |
* | --- |
* | "Rescale up" |
*
* @param {Server_Detail_Rescale_UpInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_rescale_up = /** @type {((inputs?: Server_Detail_Rescale_UpInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Rescale_UpInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_rescale_up(inputs)
});
/**
* | output |
* | --- |
* | "Rescale down" |
*
* @param {Server_Detail_Rescale_DownInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_rescale_down = /** @type {((inputs?: Server_Detail_Rescale_DownInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Rescale_DownInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_rescale_down(inputs)
});
/**
* | output |
* | --- |
* | "Request promote" |
*
* @param {Server_Detail_PromoteInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_promote = /** @type {((inputs?: Server_Detail_PromoteInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_PromoteInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_promote(inputs)
});
/**
* | output |
* | --- |
* | "Request demote" |
*
* @param {Server_Detail_DemoteInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_demote = /** @type {((inputs?: Server_Detail_DemoteInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_DemoteInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_demote(inputs)
});
/**
* | output |
* | --- |
* | "Edit windows" |
*
* @param {Server_Detail_Edit_WindowsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_edit_windows = /** @type {((inputs?: Server_Detail_Edit_WindowsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Edit_WindowsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_edit_windows(inputs)
});
/**
* | output |
* | --- |
* | "Recent events" |
*
* @param {Server_Detail_Recent_EventsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_recent_events = /** @type {((inputs?: Server_Detail_Recent_EventsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Recent_EventsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_recent_events(inputs)
});
/**
* | output |
* | --- |
* | "No events yet." |
*
* @param {Server_Detail_Events_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_events_empty = /** @type {((inputs?: Server_Detail_Events_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Events_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_events_empty(inputs)
});
/**
* | output |
* | --- |
* | "No windows yet." |
*
* @param {Server_Detail_Windows_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_windows_empty = /** @type {((inputs?: Server_Detail_Windows_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Windows_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_windows_empty(inputs)
});
/**
* | output |
* | --- |
* | "Existing windows ({count})" |
*
* @param {Server_Detail_Windows_CountInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_windows_count = /** @type {((inputs: Server_Detail_Windows_CountInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Windows_CountInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_windows_count(inputs)
});
/**
* | output |
* | --- |
* | "Enabled" |
*
* @param {Server_Detail_Window_EnabledInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_window_enabled = /** @type {((inputs?: Server_Detail_Window_EnabledInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Window_EnabledInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_window_enabled(inputs)
});
/**
* | output |
* | --- |
* | "Disabled" |
*
* @param {Server_Detail_Window_DisabledInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_window_disabled = /** @type {((inputs?: Server_Detail_Window_DisabledInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Window_DisabledInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_window_disabled(inputs)
});
/**
* | output |
* | --- |
* | "Every day" |
*
* @param {Server_Detail_Window_Every_DayInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_detail_window_every_day = /** @type {((inputs?: Server_Detail_Window_Every_DayInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Detail_Window_Every_DayInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_detail_window_every_day(inputs)
});
/**
* | output |
* | --- |
* | "Edit server" |
*
* @param {Server_Edit_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_title = /** @type {((inputs?: Server_Edit_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_title(inputs)
});
/**
* | output |
* | --- |
* | "Save" |
*
* @param {Server_Edit_SaveInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_save = /** @type {((inputs?: Server_Edit_SaveInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_SaveInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_save(inputs)
});
/**
* | output |
* | --- |
* | "Saving…" |
*
* @param {Server_Edit_SavingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_saving = /** @type {((inputs?: Server_Edit_SavingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_SavingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_saving(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Server_Edit_CancelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_cancel = /** @type {((inputs?: Server_Edit_CancelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_CancelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Server_Edit_Field_NameInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_name = /** @type {((inputs?: Server_Edit_Field_NameInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_NameInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_name(inputs)
});
/**
* | output |
* | --- |
* | "Label" |
*
* @param {Server_Edit_Field_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_label = /** @type {((inputs?: Server_Edit_Field_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_label(inputs)
});
/**
* | output |
* | --- |
* | "Base server type" |
*
* @param {Server_Edit_Field_BaseInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_base = /** @type {((inputs?: Server_Edit_Field_BaseInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_BaseInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_base(inputs)
});
/**
* | output |
* | --- |
* | "Top server type" |
*
* @param {Server_Edit_Field_TopInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_top = /** @type {((inputs?: Server_Edit_Field_TopInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_TopInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_top(inputs)
});
/**
* | output |
* | --- |
* | "Fallback chain (comma-separated, top first)" |
*
* @param {Server_Edit_Field_FallbackInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_fallback = /** @type {((inputs?: Server_Edit_Field_FallbackInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_FallbackInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_fallback(inputs)
});
/**
* | output |
* | --- |
* | "cpx31, cpx21, cpx11" |
*
* @param {Server_Edit_Field_Fallback_PlaceholderInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_fallback_placeholder = /** @type {((inputs?: Server_Edit_Field_Fallback_PlaceholderInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_Fallback_PlaceholderInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_fallback_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Mode" |
*
* @param {Server_Edit_Field_ModeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_mode = /** @type {((inputs?: Server_Edit_Field_ModeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_ModeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_mode(inputs)
});
/**
* | output |
* | --- |
* | "Timezone (IANA)" |
*
* @param {Server_Edit_Field_TimezoneInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_timezone = /** @type {((inputs?: Server_Edit_Field_TimezoneInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_TimezoneInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_timezone(inputs)
});
/**
* | output |
* | --- |
* | "Europe/Rome" |
*
* @param {Server_Edit_Field_Timezone_PlaceholderInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const server_edit_field_timezone_placeholder = /** @type {((inputs?: Server_Edit_Field_Timezone_PlaceholderInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Server_Edit_Field_Timezone_PlaceholderInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.server_edit_field_timezone_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Windows" |
*
* @param {Windows_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_title = /** @type {((inputs?: Windows_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_title(inputs)
});
/**
* | output |
* | --- |
* | "Add window" |
*
* @param {Windows_AddInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_add = /** @type {((inputs?: Windows_AddInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_AddInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_add(inputs)
});
/**
* | output |
* | --- |
* | "No windows yet." |
*
* @param {Windows_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_empty = /** @type {((inputs?: Windows_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_empty(inputs)
});
/**
* | output |
* | --- |
* | "Label" |
*
* @param {Windows_Col_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_label = /** @type {((inputs?: Windows_Col_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_label(inputs)
});
/**
* | output |
* | --- |
* | "Days" |
*
* @param {Windows_Col_DaysInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_days = /** @type {((inputs?: Windows_Col_DaysInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_DaysInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_days(inputs)
});
/**
* | output |
* | --- |
* | "Start" |
*
* @param {Windows_Col_StartInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_start = /** @type {((inputs?: Windows_Col_StartInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_StartInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_start(inputs)
});
/**
* | output |
* | --- |
* | "Stop" |
*
* @param {Windows_Col_StopInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_stop = /** @type {((inputs?: Windows_Col_StopInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_StopInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_stop(inputs)
});
/**
* | output |
* | --- |
* | "Target type" |
*
* @param {Windows_Col_TargetInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_target = /** @type {((inputs?: Windows_Col_TargetInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_TargetInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_target(inputs)
});
/**
* | output |
* | --- |
* | "Enabled" |
*
* @param {Windows_Col_EnabledInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_enabled = /** @type {((inputs?: Windows_Col_EnabledInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_EnabledInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_enabled(inputs)
});
/**
* | output |
* | --- |
* | "Yes" |
*
* @param {Windows_Col_YesInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_yes = /** @type {((inputs?: Windows_Col_YesInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_YesInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_yes(inputs)
});
/**
* | output |
* | --- |
* | "No" |
*
* @param {Windows_Col_NoInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_col_no = /** @type {((inputs?: Windows_Col_NoInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Col_NoInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_col_no(inputs)
});
/**
* | output |
* | --- |
* | "Enable" |
*
* @param {Windows_EnableInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_enable = /** @type {((inputs?: Windows_EnableInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_EnableInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_enable(inputs)
});
/**
* | output |
* | --- |
* | "Disable" |
*
* @param {Windows_DisableInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_disable = /** @type {((inputs?: Windows_DisableInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_DisableInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_disable(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Windows_DeleteInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_delete = /** @type {((inputs?: Windows_DeleteInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_DeleteInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_delete(inputs)
});
/**
* | output |
* | --- |
* | "Delete this window?" |
*
* @param {Windows_Delete_ConfirmInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_delete_confirm = /** @type {((inputs?: Windows_Delete_ConfirmInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Delete_ConfirmInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_delete_confirm(inputs)
});
/**
* | output |
* | --- |
* | "New window" |
*
* @param {Windows_Modal_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_modal_title = /** @type {((inputs?: Windows_Modal_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Modal_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_modal_title(inputs)
});
/**
* | output |
* | --- |
* | "Create" |
*
* @param {Windows_Modal_SaveInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_modal_save = /** @type {((inputs?: Windows_Modal_SaveInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Modal_SaveInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_modal_save(inputs)
});
/**
* | output |
* | --- |
* | "Creating…" |
*
* @param {Windows_Modal_SavingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_modal_saving = /** @type {((inputs?: Windows_Modal_SavingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Modal_SavingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_modal_saving(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Windows_Modal_CancelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_modal_cancel = /** @type {((inputs?: Windows_Modal_CancelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Modal_CancelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_modal_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Label" |
*
* @param {Windows_Field_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_label = /** @type {((inputs?: Windows_Field_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_label(inputs)
});
/**
* | output |
* | --- |
* | "Days of week" |
*
* @param {Windows_Field_DaysInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_days = /** @type {((inputs?: Windows_Field_DaysInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_DaysInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_days(inputs)
});
/**
* | output |
* | --- |
* | "Start time" |
*
* @param {Windows_Field_StartInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_start = /** @type {((inputs?: Windows_Field_StartInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_StartInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_start(inputs)
});
/**
* | output |
* | --- |
* | "Stop time" |
*
* @param {Windows_Field_StopInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_stop = /** @type {((inputs?: Windows_Field_StopInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_StopInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_stop(inputs)
});
/**
* | output |
* | --- |
* | "Target server type" |
*
* @param {Windows_Field_TargetInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_target = /** @type {((inputs?: Windows_Field_TargetInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_TargetInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_target(inputs)
});
/**
* | output |
* | --- |
* | "Enabled" |
*
* @param {Windows_Field_EnabledInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const windows_field_enabled = /** @type {((inputs?: Windows_Field_EnabledInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Windows_Field_EnabledInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.windows_field_enabled(inputs)
});
/**
* | output |
* | --- |
* | "Events" |
*
* @param {Events_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_title = /** @type {((inputs?: Events_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_title(inputs)
});
/**
* | output |
* | --- |
* | "Server" |
*
* @param {Events_Filter_ServerInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_filter_server = /** @type {((inputs?: Events_Filter_ServerInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_Filter_ServerInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_filter_server(inputs)
});
/**
* | output |
* | --- |
* | "All servers" |
*
* @param {Events_Filter_Server_AllInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_filter_server_all = /** @type {((inputs?: Events_Filter_Server_AllInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_Filter_Server_AllInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_filter_server_all(inputs)
});
/**
* | output |
* | --- |
* | "Kind" |
*
* @param {Events_Filter_KindInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_filter_kind = /** @type {((inputs?: Events_Filter_KindInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_Filter_KindInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_filter_kind(inputs)
});
/**
* | output |
* | --- |
* | "All kinds" |
*
* @param {Events_Filter_Kind_AllInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_filter_kind_all = /** @type {((inputs?: Events_Filter_Kind_AllInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_Filter_Kind_AllInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_filter_kind_all(inputs)
});
/**
* | output |
* | --- |
* | "Limit" |
*
* @param {Events_Filter_LimitInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_filter_limit = /** @type {((inputs?: Events_Filter_LimitInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_Filter_LimitInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_filter_limit(inputs)
});
/**
* | output |
* | --- |
* | "Loading…" |
*
* @param {Events_LoadingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_loading = /** @type {((inputs?: Events_LoadingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_LoadingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_loading(inputs)
});
/**
* | output |
* | --- |
* | "No events yet." |
*
* @param {Events_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const events_empty = /** @type {((inputs?: Events_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Events_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.events_empty(inputs)
});
/**
* | output |
* | --- |
* | "Loading…" |
*
* @param {Dashboard_LoadingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_loading = /** @type {((inputs?: Dashboard_LoadingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_LoadingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_loading(inputs)
});
/**
* | output |
* | --- |
* | "Projects ({count})" |
*
* @param {Dashboard_Section_ProjectsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_section_projects = /** @type {((inputs: Dashboard_Section_ProjectsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Section_ProjectsInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_section_projects(inputs)
});
/**
* | output |
* | --- |
* | "Servers ({count})" |
*
* @param {Dashboard_Section_ServersInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_section_servers = /** @type {((inputs: Dashboard_Section_ServersInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Section_ServersInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_section_servers(inputs)
});
/**
* | output |
* | --- |
* | "Recent events" |
*
* @param {Dashboard_Section_Recent_EventsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_section_recent_events = /** @type {((inputs?: Dashboard_Section_Recent_EventsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Section_Recent_EventsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_section_recent_events(inputs)
});
/**
* | output |
* | --- |
* | "Active servers" |
*
* @param {Kpi_Active_ServersInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_active_servers = /** @type {((inputs?: Kpi_Active_ServersInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Active_ServersInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_active_servers(inputs)
});
/**
* | output |
* | --- |
* | "Servers in auto_promote or scheduled mode" |
*
* @param {Kpi_Active_Servers_HintInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_active_servers_hint = /** @type {((inputs?: Kpi_Active_Servers_HintInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Active_Servers_HintInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_active_servers_hint(inputs)
});
/**
* | output |
* | --- |
* | "Projects" |
*
* @param {Kpi_ProjectsInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_projects = /** @type {((inputs?: Kpi_ProjectsInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_ProjectsInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_projects(inputs)
});
/**
* | output |
* | --- |
* | "Projects with a stored Hetzner API token" |
*
* @param {Kpi_Projects_HintInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_projects_hint = /** @type {((inputs?: Kpi_Projects_HintInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Projects_HintInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_projects_hint(inputs)
});
/**
* | output |
* | --- |
* | "Rescales (24h, successful)" |
*
* @param {Kpi_Rescales_24h_OkInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_rescales_24h_ok = /** @type {((inputs?: Kpi_Rescales_24h_OkInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Rescales_24h_OkInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_rescales_24h_ok(inputs)
});
/**
* | output |
* | --- |
* | "Last rescale error" |
*
* @param {Kpi_Last_ErrorInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_last_error = /** @type {((inputs?: Kpi_Last_ErrorInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Last_ErrorInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_last_error(inputs)
});
/**
* | output |
* | --- |
* | "No recent failures" |
*
* @param {Kpi_No_ErrorInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_no_error = /** @type {((inputs?: Kpi_No_ErrorInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_No_ErrorInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_no_error(inputs)
});
/**
* | output |
* | --- |
* | "Loading…" |
*
* @param {Kpi_LoadingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const kpi_loading = /** @type {((inputs?: Kpi_LoadingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_LoadingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.kpi_loading(inputs)
});
/**
* | output |
* | --- |
* | "Rescaling activity" |
*
* @param {Dashboard_Chart_ActivityInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_activity = /** @type {((inputs?: Dashboard_Chart_ActivityInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_ActivityInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_activity(inputs)
});
/**
* | output |
* | --- |
* | "Cost breakdown (€ / month)" |
*
* @param {Dashboard_Chart_CostInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_cost = /** @type {((inputs?: Dashboard_Chart_CostInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_CostInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_cost(inputs)
});
/**
* | output |
* | --- |
* | "Range" |
*
* @param {Dashboard_Chart_RangeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_range = /** @type {((inputs?: Dashboard_Chart_RangeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_RangeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_range(inputs)
});
/**
* | output |
* | --- |
* | "1d" |
*
* @param {Dashboard_Chart_Range_1dInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_range_1d = /** @type {((inputs?: Dashboard_Chart_Range_1dInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_Range_1dInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_range_1d(inputs)
});
/**
* | output |
* | --- |
* | "7d" |
*
* @param {Dashboard_Chart_Range_7dInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_range_7d = /** @type {((inputs?: Dashboard_Chart_Range_7dInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_Range_7dInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_range_7d(inputs)
});
/**
* | output |
* | --- |
* | "30d" |
*
* @param {Dashboard_Chart_Range_30dInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_range_30d = /** @type {((inputs?: Dashboard_Chart_Range_30dInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_Range_30dInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_range_30d(inputs)
});
/**
* | output |
* | --- |
* | "No rescaling history in this range." |
*
* @param {Dashboard_Chart_Cost_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_cost_empty = /** @type {((inputs?: Dashboard_Chart_Cost_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_Cost_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.dashboard_chart_cost_empty(inputs)
});
/**
* | output |
* | --- |
* | "System health" |
*
* @param {Health_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_title = /** @type {((inputs?: Health_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_title(inputs)
});
/**
* | output |
* | --- |
* | "Go API reachable" |
*
* @param {Health_Card_Api_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_api_label = /** @type {((inputs?: Health_Card_Api_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Api_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_api_label(inputs)
});
/**
* | output |
* | --- |
* | "SQLite writable" |
*
* @param {Health_Card_Db_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_db_label = /** @type {((inputs?: Health_Card_Db_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Db_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_db_label(inputs)
});
/**
* | output |
* | --- |
* | "Hetzner API reachable" |
*
* @param {Health_Card_Hcloud_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_hcloud_label = /** @type {((inputs?: Health_Card_Hcloud_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Hcloud_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_hcloud_label(inputs)
});
/**
* | output |
* | --- |
* | "Last event age" |
*
* @param {Health_Card_Last_Event_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_last_event_label = /** @type {((inputs?: Health_Card_Last_Event_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Last_Event_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_last_event_label(inputs)
});
/**
* | output |
* | --- |
* | "Recent errors (24h)" |
*
* @param {Health_Card_Recent_Errors_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_recent_errors_label = /** @type {((inputs?: Health_Card_Recent_Errors_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Recent_Errors_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_recent_errors_label(inputs)
});
/**
* | output |
* | --- |
* | "Active rescale windows" |
*
* @param {Health_Card_Windows_LabelInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_card_windows_label = /** @type {((inputs?: Health_Card_Windows_LabelInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Card_Windows_LabelInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_card_windows_label(inputs)
});
/**
* | output |
* | --- |
* | "Warning threshold: {threshold}" |
*
* @param {Health_Warn_ThresholdInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_warn_threshold = /** @type {((inputs: Health_Warn_ThresholdInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Warn_ThresholdInputs, { locale?: "en" }, {}>} */ ((inputs, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_warn_threshold(inputs)
});
/**
* | output |
* | --- |
* | "Healthy" |
*
* @param {Health_Ok_BelowInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_ok_below = /** @type {((inputs?: Health_Ok_BelowInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Ok_BelowInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_ok_below(inputs)
});
/**
* | output |
* | --- |
* | "Warning" |
*
* @param {Health_Warn_AboveInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_warn_above = /** @type {((inputs?: Health_Warn_AboveInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Warn_AboveInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_warn_above(inputs)
});
/**
* | output |
* | --- |
* | "Failing" |
*
* @param {Health_Fail_AboveInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_fail_above = /** @type {((inputs?: Health_Fail_AboveInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_Fail_AboveInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_fail_above(inputs)
});
/**
* | output |
* | --- |
* | "Refresh" |
*
* @param {Health_CheckingInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const health_checking = /** @type {((inputs?: Health_CheckingInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Health_CheckingInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.health_checking(inputs)
});
/**
* | output |
* | --- |
* | "Server status" |
*
* @param {Servers_Status_TitleInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_title = /** @type {((inputs?: Servers_Status_TitleInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_TitleInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_title(inputs)
});
/**
* | output |
* | --- |
* | "No servers registered." |
*
* @param {Servers_Status_EmptyInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_empty = /** @type {((inputs?: Servers_Status_EmptyInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_EmptyInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_empty(inputs)
});
/**
* | output |
* | --- |
* | "Server" |
*
* @param {Servers_Status_Col_NameInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_col_name = /** @type {((inputs?: Servers_Status_Col_NameInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_Col_NameInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_col_name(inputs)
});
/**
* | output |
* | --- |
* | "Mode" |
*
* @param {Servers_Status_Col_ModeInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_col_mode = /** @type {((inputs?: Servers_Status_Col_ModeInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_Col_ModeInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_col_mode(inputs)
});
/**
* | output |
* | --- |
* | "Top type" |
*
* @param {Servers_Status_Col_TopInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_col_top = /** @type {((inputs?: Servers_Status_Col_TopInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_Col_TopInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_col_top(inputs)
});
/**
* | output |
* | --- |
* | "Active window" |
*
* @param {Servers_Status_Col_WindowInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_col_window = /** @type {((inputs?: Servers_Status_Col_WindowInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_Col_WindowInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_col_window(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Servers_Status_Col_StatusInputs} inputs
* @param {{ locale?: "en" }} options
* @returns {LocalizedString}
*/
export const servers_status_col_status = /** @type {((inputs?: Servers_Status_Col_StatusInputs, options?: { locale?: "en" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Servers_Status_Col_StatusInputs, { locale?: "en" }, {}>} */ ((inputs = {}, options = {}) => {
	experimentalStaticLocale ?? options.locale ?? getLocale()
	return __en.servers_status_col_status(inputs)
});