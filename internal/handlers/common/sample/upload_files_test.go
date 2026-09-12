package sample_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createFormFile(field, file string) (*bytes.Buffer, *multipart.Writer) {
	return createFormFileContent(field, file, "dummy")
}

func TestUploadFiles(t *testing.T) {
	testutils.SetupTestContext()
	mockUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")
		dir := t.TempDir()

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (
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

		expectedFilePath := filepath.Join(dir, "reads_R1.fastq")
		fileContent, err := os.ReadFile(expectedFilePath)

		assert.NoError(t, err)
		assert.Equal(t, "dummy", string(fileContent))

		expected := testutils.ToJSON(map[string]string{
			"message": "Sample files submitted successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Collaborator cannot upload to another user's sample",
		func(t *testing.T) {
			buf, mw := createFormFile("fasta", "contigs.fasta")
			dir := t.TempDir()

			svc := &mocks.MockSampleService{
				GetSampleForUploadFunc: func(_ context.Context,
					_ uuid.UUID) (*models.Sample, error) {
					sample := testmodels.CreateMockSample()
					sample.UserID = uuid.New()
					return &sample, nil
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
				"error": "Unauthorized. Please log in to continue.",
			})

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Equal(t, expected, w.Body.String())

			entries, err := os.ReadDir(dir)

			assert.NoError(t, err)
			assert.Empty(t, entries, "no file should be written")
		})

	t.Run("Error - GetSampleForUpload Not Found", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				return nil, services.ErrNotFound
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

	t.Run("Error - GetSampleForUpload Internal Error", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				return nil, services.ErrInternal
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
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")
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
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (
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
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (
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
		buf, mw := createFormFile("fastq2", "reads_R2.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
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
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
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
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
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
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (
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
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")

		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (
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

func createFormFileContent(field, file, content string) (*bytes.Buffer,
	*multipart.Writer) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fw, _ := mw.CreateFormFile(field, file)
	io.WriteString(fw, content)

	mw.Close()
	return &buf, mw
}

func TestUploadFilesValidations(t *testing.T) {
	testutils.SetupTestContext()
	mockUserID := uuid.New()

	t.Run("Success - Unknown Form Field Ignored", func(t *testing.T) {
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)

		fw, _ := mw.CreateFormFile("evil", "virus.exe")
		io.WriteString(fw, "dummy")
		fw, _ = mw.CreateFormFile("fasta", "contigs.fasta")
		io.WriteString(fw, "dummy")
		mw.Close()

		dir := t.TempDir()
		svc := &mocks.MockSampleService{
			GetSampleForUploadFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Sample, error) {
				sample := testmodels.CreateMockSample()
				sample.UserID = mockUserID
				return &sample, nil
			},
			PrepareSampleFolderFunc: func(_ context.Context, userID,
				sampleID uuid.UUID) (string, error) {
				return dir, nil
			},
			AttachFilesFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleAttachmentInput) error {
				assert.NotNil(t, input.Fasta)
				assert.Nil(t, input.Fastq1)
				return nil
			},
		}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", &buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.FileExists(t, filepath.Join(dir, "contigs.fasta"))
		assert.NoFileExists(t, filepath.Join(dir, "virus.exe"))
	})

	t.Run("Error - Unsupported Extension", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "virus.exe")
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "File type not supported.")
		assert.NoFileExists(t, filepath.Join(dir, "virus.exe"))
	})

	t.Run("Error - Gzipped FASTA Rejected", func(t *testing.T) {
		buf, mw := createFormFile("fasta", "contigs.fasta.gz")
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "File type not supported.")
		assert.NoFileExists(t, filepath.Join(dir, "contigs.fasta.gz"))
	})

	t.Run("Error - File Name Too Long", func(t *testing.T) {
		fileName := strings.Repeat("a", 260) + ".fastq"
		buf, mw := createFormFile("fastq1", fileName)
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "File type not supported.")
	})

	t.Run("Error - Invalid Gzip Magic", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq.gz")
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "File type not supported.")
		assert.NoFileExists(t, filepath.Join(dir, "reads_R1.fastq.gz"))
	})

	t.Run("Success - Valid Gzip Magic", func(t *testing.T) {
		content := "\x1f\x8b" + "junk"
		buf, mw := createFormFileContent("fastq1", "reads_R1.fastq.gz",
			content)
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		fileContent, err := os.ReadFile(filepath.Join(dir, "reads_R1.fastq.gz"))

		assert.NoError(t, err)
		assert.Equal(t, content, string(fileContent))
	})

	t.Run("Error - Upload Exceeds Budget", func(t *testing.T) {
		original := config.MaxUploadSize
		config.MaxUploadSize = 1
		defer func() { config.MaxUploadSize = original }()

		buf, mw := createFormFile("fastq1", "reads_R1.fastq")
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, nil)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "The file exceeds the maximum upload size.")
		assert.NoFileExists(t, filepath.Join(dir, "reads_R1.fastq"))
	})

	t.Run("Error - AttachFiles Failure Removes Files", func(t *testing.T) {
		buf, mw := createFormFile("fastq1", "reads_R1.fastq")
		dir := t.TempDir()

		svc := uploadMock(mockUserID, dir, services.ErrMissingFiles)
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinMultipartContext(
			http.MethodPut, "/api/sample", buf, mw.FormDataContentType(),
			nil, gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.UploadFiles(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		entries, err := os.ReadDir(dir)

		assert.NoError(t, err)
		assert.Empty(t, entries, "saved files should be removed on failure")
	})
}

func uploadMock(userID uuid.UUID, dir string,
	attachErr error) *mocks.MockSampleService {
	return &mocks.MockSampleService{
		GetSampleForUploadFunc: func(_ context.Context,
			_ uuid.UUID) (*models.Sample, error) {
			sample := testmodels.CreateMockSample()
			sample.UserID = userID
			return &sample, nil
		},
		PrepareSampleFolderFunc: func(_ context.Context, _,
			_ uuid.UUID) (string, error) {
			return dir, nil
		},
		AttachFilesFunc: func(ctx context.Context, sampleID,
			uid uuid.UUID, input models.SampleAttachmentInput) error {
			return attachErr
		},
	}
}
