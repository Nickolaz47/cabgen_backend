package sample_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createFormFile(field, file string) (*bytes.Buffer, *multipart.Writer) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fw, _ := mw.CreateFormFile(field, file)
	io.WriteString(fw, "dummy")

	mw.Close()
	return &buf, mw
}

func TestUploadFiles(t *testing.T) {
	testutils.SetupTestContext()
	mockUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")
		dir := t.TempDir()

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID, sampleID uuid.UUID) (
				string, error) {

				assert.Equal(t, mockUserID, userID)
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {

				assert.Equal(t, mockUserID, userID)
				return nil
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expectedFilePath := filepath.Join(dir, "reads_R1.fastq.gz")
		fileContent, err := os.ReadFile(expectedFilePath)

		assert.NoError(t, err)
		assert.Equal(t, "dummy", string(fileContent))

		expected := testutils.ToJSON(map[string]string{
			"message": "Sample files submitted successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "The URL ID is invalid."})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID, sampleID uuid.UUID) (
				string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrNotFound
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Sample not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID, sampleID uuid.UUID) (
				string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrNotFound
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)

		handler.UploadFiles(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Unauthorized. Please log in to continue.",
		})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Content-Type", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrNotFound
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			"",
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The request must be multipart/form-data.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - AttachFiles Missing Fastq1", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq2", "reads_R2.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID,
				sampleID uuid.UUID) (string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrMissingFastq1
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The Fastq1 file was not sent.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - AttachFiles Missing Fastq2", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID,
				sampleID uuid.UUID) (string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrMissingFastq2
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The Fastq2 file was not sent.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - AttachFiles Missing Files", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("", "")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID,
				sampleID uuid.UUID) (string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrMissingFiles
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No files were sent for upload.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - PrepareSampleFolder Internal Error", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID, sampleID uuid.UUID) (
				string, error) {
				return dir, services.ErrInternal
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - AttachFiles Internal Error", func(t *testing.T) {
		dir := t.TempDir()
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(userID, sampleID uuid.UUID) (
				string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				return services.ErrInternal
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut,
			"/api/sample",
			buf,
			mw.FormDataContentType(),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})
}
