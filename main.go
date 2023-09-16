package main

import (
	"crypto/rand"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/foomo/htpasswd"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Hàm tạo mật khẩu giả mạo với độ dài cho trước
func generateFakePassword(length int) string {
	randomPassword := make([]byte, length)
	for i := range randomPassword {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		randomPassword[i] = charset[index.Int64()]
	}
	return string(randomPassword)
}
func containsString(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func generateRandomIPv6() string {
	ip := make([]byte, 16)

	// Đặt các giá trị cố định để đảm bảo định dạng đúng của IPv6
	ip[0] = 0x20
	ip[1] = 0x01

	// Tạo số ngẫu nhiên cho 12 byte cuối
	_, err := rand.Read(ip[8:])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("2001:19f0:7001:321a:%02x:%02x:%02x:%02x", ip[8], ip[9], ip[10], ip[11])
}

func main() {
	// Đường dẫn tới file cần tìm dòng chứa "x.x.x.x"
	filePath := "./config24.conf"
	fileExcel, _ := excelize.OpenFile("./data1.xlsx")
	filePassword := "./pass.htpasswd"
	// Mở file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Lỗi khi mở file:", err)
		return
	}
	defer file.Close()

	// Đọc nội dung file vào một biến
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Lỗi khi đọc file:", err)
		return
	}

	// Chuyển nội dung file thành chuỗi
	fileString := string(content)

	//index := strings.Index(fileString, "#aclport")
	//if index == -1 {
	//	fmt.Println("Không tìm thấy dòng chứa '#aclport' trong file.")
	//	return
	//}
	//aclport := []string{}
	//for i := 1; i <= 300; i++ {
	//	aclport = append(aclport, fmt.Sprintf("acl port%d localport %d", i, i+24000))
	//}
	//contenAfterChangeAcl := fileString[:index+len("#aclport")] + "\n" + strings.Join(aclport, "\n") + "\n" + fileString[index+len("#aclport"):]

	index := strings.Index(fileString, "#userauthenproxy")
	if index == -1 {
		fmt.Println("Không tìm thấy dòng chứa '#userauthenproxy' trong file.")
		return
	}
	userauthen := []string{}
	userNames := []string{}
	for i := 1; i <= 300; i++ {
		username := ""
		for {
			username = strings.ToLower(gofakeit.Username())
			if !containsString(userNames, username) {
				break
			}
		}
		userNames = append(userNames, username)
		fakePassword := generateFakePassword(10)
		_ = htpasswd.SetPassword(filePassword, username, fakePassword, htpasswd.HashAPR1)

		fileExcel.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), username)
		fileExcel.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), fakePassword)
		userauthen = append(userauthen, fmt.Sprintf("acl user%d proxy_auth %s", i, username))
	}
	contenAfterChangeAcl := fileString[:index+len("#userauthenproxy")] + "\n" + strings.Join(userauthen, "\n") + "\n" + fileString[index+len("#userauthenproxy"):]

	//index := strings.Index(fileString, "#httpaccess")
	//if index == -1 {
	//	fmt.Println("Không tìm thấy dòng chứa '#httpaccess' trong file.")
	//	return
	//}
	//aclport := []string{}
	//for i := 1; i <= 300; i++ {
	//	aclport = append(aclport, fmt.Sprintf("http_access allow user%d port%d", i, i))
	//}
	//contenAfterChangeAcl := fileString[:index+len("#httpaccess")] + "\n" + strings.Join(aclport, "\n") + "\n" + fileString[index+len("#httpaccess"):]

	//index := strings.Index(fileString, "#httpport")
	//if index == -1 {
	//	fmt.Println("Không tìm thấy dòng chứa '#httpport' trong file.")
	//	return
	//}
	//aclport := []string{}
	//for i := 1; i <= 300; i++ {
	//	aclport = append(aclport, fmt.Sprintf("http_port 0.0.0.0:%d", i+24000))
	//}
	//contenAfterChangeAcl := fileString[:index+len("#httpport")] + "\n" + strings.Join(aclport, "\n") + "\n" + fileString[index+len("#httpport"):]

	//index := strings.Index(fileString, "#tcpoutgoingaddress")
	//if index == -1 {
	//	fmt.Println("Không tìm thấy dòng chứa '#httpport' trong file.")
	//	return
	//}
	//aclport := []string{}
	//for i := 1; i <= 300; i++ {
	//	ipv6 := generateRandomIPv6()
	//	fileExcel.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), ipv6)
	//	aclport = append(aclport, fmt.Sprintf("tcp_outgoing_address %s port%d\ntcp_outgoing_address 127.0.0.1 port%d", ipv6, i, i))
	//}
	//contenAfterChangeAcl := fileString[:index+len("#tcpoutgoingaddress")] + "\n" + strings.Join(aclport, "\n") + "\n" + fileString[index+len("#tcpoutgoingaddress"):]

	err = ioutil.WriteFile(filePath, []byte(contenAfterChangeAcl), 0644)
	if err != nil {
		fmt.Println("Lỗi khi ghi file:", err)
		return
	}

	fmt.Println("Đã thay đổi thành công!")

	if err := fileExcel.Save(); err != nil {
		fmt.Println(err)
	}
}
