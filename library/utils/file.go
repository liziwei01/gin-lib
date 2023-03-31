/*
 * @Author: liziwei01
 * @Date: 2022-03-20 17:31:44
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-20 17:32:49
 * @Description: file content
 */
package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

const BaseDir = "./temp_file/"

/**
 * @description: 根据FileHeader获取文件二进制数组
 * @param {*multipart.FileHeader} fileHeader
 * @return {*}
 */
func (u *UFile) GetFileBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	if fileHeader == nil {
		return nil, fmt.Errorf("[GetFileBytes]: empty fileHeader")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func (u *UFile) IsExist(f string) bool {
	_, err := os.Stat(BaseDir + f)
	return err == nil || os.IsExist(err)
}

//CreateDir  文件夹创建
func (u *UFile) CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// SaveFile 保存文件到本地临时文件夹内
func (u *UFile) SaveFile(file multipart.File, filename string) (err error) {
	if !u.IsExist(BaseDir) {
		err := u.CreateDir(BaseDir)
		if err != nil {
			return err
		}
	}
	newFile, err := os.Create(BaseDir + filename)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}
	return nil
}

// DelFile 删除临时文件夹内的文件
func (u *UFile) DelFile(filename string) (err error) {
	if strings.Trim(filename, " ") == "" {
		return errors.New("file path is empty")
	}
	er := os.Remove(BaseDir + filename)
	if er != nil {
		return er
	}
	return nil
}

// 获取临时文件夹内的临时文件路径
func (u *UFile) GetFilePath(filename string) string {
	return BaseDir + filename
}

func (u *UFile) ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(BaseDir + filename)
	if err != nil {
		return nil, err
	}
	byteStream, err := ioutil.ReadAll(file)
	return byteStream, err
}
