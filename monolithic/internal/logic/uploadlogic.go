package logic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"monolithic/internal/svc"
	"monolithic/internal/types"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20 // 10 MB

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadLogic(r *http.Request, svcCtx *svc.ServiceContext) UploadLogic {
	return UploadLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload() (resp *types.UploadResponse, err error) {
	l.r.ParseMultipartForm(maxFileSize)
	file, handler, err := l.r.FormFile("myFile")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return &types.UploadResponse{
			Code: 1,
		}, nil
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("udmuser", "B2")
	if err != nil {
		return
	}
	logx.Infof("Sheet1 B2: %+v", cell)
	logx.Infof("upload file: %+v, file size: %d, MIME header: %+v",
		handler.Filename, handler.Size, handler.Header)

	tempFile, err := os.Create(path.Join(l.svcCtx.Config.Path, handler.Filename))
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	return &types.UploadResponse{
		Code: 0,
	}, nil
}
