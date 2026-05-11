// Wails v3 AuthService bridge
// Wraps generated window.go.main.AuthService calls as ES module exports

const auth = () => window.go.main.AuthService

export function IsLoggedIn() { return auth().IsLoggedIn() }
export function GetUserProfile() { return auth().GetUserProfile() }
export function CreateQrcode() { return auth().CreateQrcode() }
export function CheckQrcodeStatus(key) { return auth().CheckQrcodeStatus(key) }
export function SendSmsCode(phone, countryCode) { return auth().SendSmsCode(phone, countryCode) }
export function LoginByPhone(phone, password, captcha, countryCode) { return auth().LoginByPhone(phone, password, captcha, countryCode) }
export function LoginByCookie(cookieStr) { return auth().LoginByCookie(cookieStr) }
export function LoginByCookieFile(filePath) { return auth().LoginByCookieFile(filePath) }
export function LoginByCookieCloud(server, uuid, password) { return auth().LoginByCookieCloud(server, uuid, password) }
export function Logout() { return auth().Logout() }
