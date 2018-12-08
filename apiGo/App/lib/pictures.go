package lib

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func generatePng(path string, res io.Reader) (string, error) {
	img, err := png.Decode(res)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(path+".png", os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return "", err
	}
	png.Encode(f, img)
	return path + ".png", nil
}

func generateJpeg(path string, res io.Reader) (string, error) {
	img, err := jpeg.Decode(res)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(path+".jpeg", os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return "", err
	}
	jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	return path + ".jpeg", nil
}

func GeneratePicture(typeImage, pathWithoutExtension string, pictureReadable io.Reader) (string, error) {
	var imagePath string
	var err error
	switch typeImage {
	case "image/png":
		imagePath, err = generatePng(pathWithoutExtension, pictureReadable)
		if err != nil {
			return "", errors.New("Corrupted file [" + typeImage + "] | Error: " + err.Error())
		}
	case "image/jpg":
		imagePath, err = generateJpeg(pathWithoutExtension, pictureReadable)
		if err != nil {
			return "", errors.New("Corrupted file [" + typeImage + "] | Error: " + err.Error())
		}
	case "image/jpeg":
		imagePath, err = generateJpeg(pathWithoutExtension, pictureReadable)
		if err != nil {
			return "", errors.New("Corrupted file [" + typeImage + "] | Error: " + err.Error())
		}
	default:
		return "", errors.New("Image type [" + typeImage + "] not accepted, support only png, jpg and jpeg images")
	}
	return imagePath, nil
}

func ExtractDataPictureBase64(base64 string) (string, string, error) {
	if base64 == "" {
		return "", "", errors.New("Base64 can't be empty")
	}
	if !IsValidBase64Picture(base64) {
		return "", "", errors.New("Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'")
	}
	preBase64 := TrimStringFromString(base64, ";base64")
	typeImage := string(preBase64)[5:]
	imageBase64 := string(base64)[len(preBase64)+8:]
	return typeImage, imageBase64, nil
}

func Base64ToImageFile(imageBase64, typeImage, path, subPath, username string) (string, error) {
	newpath := path + subPath + username
	os.MkdirAll(newpath, os.ModePerm)
	fileName := GetRandomString(43)
	unbased, _ := Base64Decode(imageBase64) // No need to check the error here, this will be handled just after
	pictureReadable := bytes.NewReader(unbased)
	imagePath, err := GeneratePicture(typeImage, newpath+"/"+fileName, pictureReadable)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(imagePath, path), nil
}
