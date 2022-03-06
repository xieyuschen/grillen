package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Pivot-Studio/HUSTHoleBackEnd/consts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//CheckError returns false if err is not nil and output a info as log
func CheckError(err error) bool {
	if err != nil {
		logrus.Error(err)
		return true
	}
	return false
}

func IsEmailAllowedRegister(email string) (IsPermitted bool) {
	email = strings.ReplaceAll(email, " ", "")
	return strings.HasSuffix(email, "@hust.edu.cn")
}

//CheckHustEmailSuffixOrAddIt 为邮箱添加hust.edu.cn的后缀如果已有则不变，同时起检查邮箱的作用
func CheckHustEmailSuffixOrAddIt(EmailOrUid string) (email string) {
	emailSuffix := "@hust.edu.cn"
	if !strings.Contains(EmailOrUid, emailSuffix) {
		email = EmailOrUid + emailSuffix
	} else {
		email = EmailOrUid
	}
	return email
}

/*
IncrementallyUpdate 增量修改结构体变量

modifiedValue:需要被修改的结构体，其值会被info中非空的值所取代

info:提供修改的内容，info中所有非空的值都会被赋给modifiedValue

example:


type Demo struct{
	id    		int
	name  		string
	description string
}
func main(){
	update:=Demo{1,"hello","This is origin demo"}
	info :=Demo{description:"I want to change it")
	IncrementallyUpdate(&update,&info)

	//`update` value is {1,"hello","I want to change it"}
	//when IncrementallyUpdate is done
}
*/
func IncrementallyUpdate(modifiedValue interface{}, info interface{}) {
	updateS := reflect.ValueOf(modifiedValue).Elem()
	infoS := reflect.ValueOf(info).Elem()
	for i := 0; i < updateS.NumField(); i++ {
		f := updateS.Field(i)
		fInfo := infoS.Field(i)
		if !fInfo.IsZero() {
			f.Set(fInfo)
		}
	}
}

//GenerateTokenByJwt 使用Jwt生成token作身份认证使用
//其中存储了登录的时间，登陆的email以及用户的角色信息
func GenerateTokenByJwt(email string, role string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"email":     email,
		"timeStamp": GetTimeStamp(),
		"role":      role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(consts.TOKEN_SCRECT_KEY))
	return
}

//GenerateVerifyCode 随机生成生成length长度的激活码
func GenerateVerifyCode(length int) (verifyString string) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	verifyString = string(b)
	return
}

//GenerateAliasByReplyIndexAndHoleId 根据postAliasIndex和树洞号holeId为用户随机生成昵称
func GenerateAliasByReplyIndexAndHoleId(postAliasIndex uint, holeId int) (alias string) {

	if postAliasIndex == 0 {
		alias = consts.NAME[0]
		return alias
	}
	if postAliasIndex >= consts.MAX_BASE_ALIAS_NUM {
		if postAliasIndex < consts.MAX_BASE_ALIAS_NUM+9 {
			num := strconv.Itoa(int(postAliasIndex) - consts.MAX_BASE_ALIAS_NUM + 1)
			return "520幸运宝" + num + "号"
		} else {
			return "520超级大幸运宝！"
		}
	}
	//index := consts.HASH_HOLEID_FACTOR*hole_id + consts.HASH_REPLYINDEX_FACTOR*int(post_alias_index)
	name_length := len(consts.NAME)
	start_index := holeId * holeId % name_length
	//三种情况，依据昵称库进行生成
	//如果目前的postAliasIndex小于昵称总数，即当前发言的人均可被分配一个昵称，所以直接做一个简单的hash后产生新昵称
	//若postAliasIndex大于当前昵称数，小于昵称总数的平方，将两个昵称进行拼接
	//postAliasIndex大于昵称数的平方，则随即昵称后加hash字符的后五位
	if int(postAliasIndex) < name_length {
		alias_index := (start_index+int(postAliasIndex))%(name_length-1) + 1
		alias = consts.NAME[alias_index]
	} else if int(postAliasIndex)-name_length < (name_length-1)*(name_length-1) {
		//start_index % ((name_length - 1) * (name_length - 1))
		alias_index := (start_index+int(postAliasIndex))%((name_length-1)*(name_length-1)) + 1

		name1 := consts.NAME[(alias_index / (name_length - 1))]

		//from 1 to name_length-1
		name2 := consts.NAME[alias_index%(name_length-1)+1]
		alias = name1 + "与" + name2

	} else {
		alias_index := (start_index+int(postAliasIndex))%(name_length-1) + 1
		name := consts.NAME[alias_index]
		hashStr := HashWithSalt(strconv.Itoa(alias_index))
		hashSuffix := hashStr[len(hashStr)-5:]
		alias = name + hashSuffix
	}
	return
}

/*
	AES加密函数，将明文stringToEncrypt加密成为密文encryptedString

	AES的密钥从配置文件中读入
*/
func EncryptWithAes(stringToEncrypt string) (encryptedString string, err error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return fmt.Sprintf("%x", ciphertext), nil
}

func DecryptWithAes(CipherText string) (decryptedString string, err error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}
	ciphertext, _ := hex.DecodeString(CipherText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), nil
}

//RsaDecryption chiper is a string with base64 encoding,and return a plaintext
func RsaDecryption(chiper string) (plaintext string, err error) {
	des, err := base64.StdEncoding.DecodeString(chiper)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	privateKeyStr, err := ioutil.ReadFile("private.pem")
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	priblock, _ := pem.Decode(privateKeyStr)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	priKey, err := x509.ParsePKCS1PrivateKey(priblock.Bytes)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	plain, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, des)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if err != nil {
		logrus.Error(err)
		log.Panic(err)
	}
	plaintext = string(plain)
	return plaintext, nil
}

//HashWithSalt A Hash function using salt with bcrypt libriry to hash password
//将纯字符串hash
func HashWithSalt(plainText string) (HashText string) {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.MinCost)
	CheckError(err)
	HashText = string(hash)
	return
}

//ParseInt64ToTimeType 将int64的时间格式转换为time.Time
func ParseInt64ToTimeType(timestamp int64) time.Time {
	i, err := strconv.ParseInt(strconv.Itoa(int(timestamp)), 10, 64)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tm := time.Unix(i, 0).In(loc)
	return tm
}

//GetEmailFromAuthorization Get encrypted user email
//err is not nil when authorization token is not valid
func GetEmailFromAuthorization(c *gin.Context) (email string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("您未登录，请登陆后查看")
		}
	}()
	authHeader := c.Request.Header.Get("Authorization")
	token := strings.Fields(authHeader)[1]
	if err != nil {
		err = errors.New("您未登录，请登陆后查看")
		return
	}
	claim, _ := GetClaimFromToken(token)
	email = claim.(jwt.MapClaims)["email"].(string)
	return
}

func GetTokenFromAuthorization(c *gin.Context) (token string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("您未登录，请登陆后查看")
		}
	}()
	authHeader := c.Request.Header.Get("Authorization")
	token = strings.Fields(authHeader)[1]
	if err != nil {
		err = errors.New("您未登录，请登陆后查看")
		return
	}
	return
}

func GetClaimFromToken(tokenString string) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SCRECT_KEY), err
	})
	if err != nil {
		return nil, err
	} else {
		claims = token.Claims.(jwt.MapClaims)
		return claims, nil
	}
}

func GetTimeStamp() (t int64) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = time.Now().In(loc).Unix()
	return
}

func GetTodayStartTimeStamp() (timestamp int64) {

	t := time.Now()
	year, month, day := t.Date()
	timestamp = time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Unix()
	return
}

//SplitTimeFromMillionSecondToSecond 如果是13为的时间戳（单位毫秒），那就切换成10为以毫秒为单位的
//如果为null。就直接传空字符串回去
func SplitTimeFromMillionSecondToSecond(millionSecond string) (t string) {
	if len(millionSecond) == 13 {
		t = millionSecond[:10]
	}
	return t
}

func IsReservedHoleID(holeID int) bool {
	for _, v := range consts.RESERVED_HOLE_ID {
		if holeID == v {
			return true
		}
	}
	return false
}