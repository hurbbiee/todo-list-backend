package enum

import "errors"

type ErrorDef struct {
	Code    int
	Message string
}

var (
	BadRequest = ErrorDef{
		Code:    1000,
		Message: "ข้อมูลที่ส่งมาไม่ถูกต้อง",
	}

	ValidationError = ErrorDef{
		Code:    1001,
		Message: "ข้อมูลไม่ผ่านการตรวจสอบ",
	}

	ValidationFormatDate = ErrorDef{
		Code:    1002,
		Message: "ข้อมูลวันที่ส่งมาไม่ถูกต้อง",
	}
	ValidateDateIsPast = ErrorDef{
		Code:    1003,
		Message: "ข้อมูลวันที่ส่งย้อนหลังไม่ได้",
	}

	Internal = ErrorDef{
		Code:    5000,
		Message: "เกิดข้อผิดพลาดภายในระบบ",
	}
)

var (
	UserNotFound = ErrorDef{
		Code:    2001,
		Message: "ไม่พบผู้ใช้งาน",
	}

	EmailAlreadyExists = ErrorDef{
		Code:    23505,
		Message: "อีเมลนี้มีผู้ใช้งานในระบบแล้ว",
	}
)

var (
	AuthInvalidCredential = ErrorDef{
		Code:    4011,
		Message: "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
	}

	AuthTokenExpired = ErrorDef{
		Code:    4012,
		Message: "โทเค็นหมดอายุ",
	}

	AuthUnauthorized = ErrorDef{
		Code:    4013,
		Message: "ยังไม่ได้เข้าสู่ระบบ",
	}

	PermissionDenied = ErrorDef{
		Code:    4031,
		Message: "ไม่มีสิทธิ์เข้าถึง",
	}
)

var (
	DiscordInvalidSettings = ErrorDef{
		Code:    3001,
		Message: "การตั้งค่า Discord ไม่ถูกต้อง",
	}

	DiscordInvalidWebhookURL = ErrorDef{
		Code:    3002,
		Message: "Discord Webhook URL ไม่ถูกต้อง",
	}

	DiscordConnectionNotFound = ErrorDef{
		Code:    3003,
		Message: "ไม่พบการเชื่อมต่อ Discord",
	}
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
