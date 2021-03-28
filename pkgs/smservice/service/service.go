/*
 * Revision History:
 *     Initial: 2019/03/20        Yang ChengKai
 */

package services

import (
	ran "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"time"

	"github.com/abserari/shower/smservice/model/mysql"
)

var numbers = []byte("012345678998765431234567890987654321")

// SMSVerify -
type SMSVerify interface {
	OnVerifySucceed(targetID, mobile string)
	OnVerifyFailed(targetID, mobile string)
}

// Config -
type Config struct {
	Host           string
	Appcode        string
	Digits         int
	ResendInterval int
	OnCheck        SMSVerify
}

// Controller -
type Controller struct {
	DB   *sql.DB
	Conf Config
}

// NewController -
func NewController(db *sql.DB, Conf *Config) *Controller {
	sm := &Controller{
		DB: db,
		Conf: Config{
			Host:           Conf.Host,
			Appcode:        Conf.Appcode,
			Digits:         Conf.Digits,
			ResendInterval: Conf.ResendInterval,
			OnCheck:        Conf.OnCheck,
		},
	}

	return sm
}

// SendSmsReply -
type SendSmsReply struct {
	Message   string `json:"Message,omitempty"`
	RequestID string `json:"RequestId,omitempty"`
	BizID     string `json:"BizId,omitempty"`
	Code      string `json:"Code,omitempty"`
}

// SMS -
type SMS struct {
	Mobile string
	Date   int64
	Code   string
	Sign   string
}

func newSms() *SMS {
	sms := &SMS{}
	return sms
}

//准备发送的结构
func (sms *SMS) prepare(mobile, sign string, digits int) {
	sms.Mobile = mobile
	sms.Date = time.Now().Unix()
	sms.Code = Code(digits)
	sms.Sign = sign
}

// Code 生成x位数字验证码
func Code(size int) string {
	data := make([]byte, size)
	out := make([]byte, size)

	buffer := len(numbers)
	_, err := ran.Read(data)
	if err != nil {
		panic(err)
	}

	for id, key := range data {
		x := byte(int(key) % buffer)
		out[id] = numbers[x]
	}

	return string(out)
}

//有效检验
func (sms *SMS) checkvalid(db *sql.DB, conf *Config) error {
	unixtime := sms.getDate(db)

	if unixtime > 0 && sms.Date-unixtime < int64(conf.ResendInterval) {
		return errors.New("短时间内不允许发送两次")
	}

	if err := VailMobile(sms.Mobile); err != nil {
		return errors.New("手机号不符合规则")
	}

	return nil
}

func (sms *SMS) getDate(db *sql.DB) int64 {
	unixtime, _ := mysql.GetDate(db, sms.Sign)
	return unixtime
}

// VailMobile 可行的手机号
func VailMobile(mobile string) error {
	if len(mobile) < 11 {
		return errors.New("[mobile]参数不对")
	}

	reg, err := regexp.Compile("^1[3-8][0-9]{9}$")
	if err != nil {
		panic("regexp error")
	}

	if !reg.MatchString(mobile) {
		return errors.New("手机号码[mobile]格式不正确")
	}

	return nil
}

//发送后存储这个信息：手机号,时间，验证码
//存储入数据库
func (sms *SMS) save(db *sql.DB) error {
	err := mysql.Insert(db, sms.Mobile, sms.Date, sms.Code, sms.Sign)

	return err
}

//实现一个可以直接改配置就能用的send方法
//1.拼接函数，拼接成需要的url
//2.设置参数
//aliyun
func (sms *SMS) sendmsg(conf *Config) error {
	host := conf.Host

	url := host + "?" + "code=" + sms.Code + "&phone=" + sms.Mobile + "&skin=1"

	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "APPCODE "+conf.Appcode)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	ssr := &SendSmsReply{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return err
	}

	if ssr.Code != "OK" {
		return errors.New(ssr.Code)
	}

	return nil
}

//Send 根据手机号和id生成时间和验证码，并发送后存入数据库
func Send(mobile, sign string, conf *Config, db *sql.DB) error {
	sms := newSms()
	sms.prepare(mobile, sign, conf.Digits)

	if err := sms.checkvalid(db, conf); err != nil {
		return err
	}

	if err := sms.save(db); err != nil {
		return err
	}

	err := sms.sendmsg(conf)

	return err
}

//Check 根据sign和验证码，返回nil表示成功
func Check(code, sign string, conf *Config, db *sql.DB) error {
	sms := newSms()
	sms.Date = time.Now().Unix()
	sms.Code = code
	sms.Sign = sign
	//验证超时

	//验证
	getcode, err := sms.getCode(db)
	if err != nil {
		return errors.New("Sign error")
	}

	if sms.Code == getcode {
		sms.delete(sms.Sign, db)
		return nil
	}

	return errors.New("Code error")
}

func (sms *SMS) getCode(db *sql.DB) (string, error) {
	code, err := mysql.GetCode(db, sms.Sign)
	return code, err
}

//删除数据库数据
func (sms *SMS) delete(sign string, db *sql.DB) { mysql.Delete(db, sign) }

// UID 生成uid
func UID() string {
	data := make([]byte, 16)

	_, err := ran.Read(data)
	if err != nil {
		panic(err)
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", data[0:4], data[4:6], data[6:8], data[8:10], data[10:])
	return uuid
}
