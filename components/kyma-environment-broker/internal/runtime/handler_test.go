package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage/driver/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeHandler(t *testing.T) {

	t.Run("test pagination should work", func(t *testing.T) {
		// given
		operations := memory.NewOperation()
		instances := memory.NewInstance(operations)
		testID1 := "Test1"
		testID2 := "Test2"

		err := instances.Insert(internal.Instance{
			InstanceID: testID1,
		})
		require.NoError(t, err)
		err = instances.Insert(internal.Instance{
			InstanceID: testID2,
		})
		require.NoError(t, err)

		runtimeHandler := NewHandler(instances, operations, 2, NewConverter())

		req, err := http.NewRequest("GET", "/runtimes?limit=1", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		runtimeHandler.AttachRoutes(router)

		// when
		router.ServeHTTP(rr, req)

		// then
		require.Equal(t, http.StatusOK, rr.Code)

		var out RuntimesPage

		err = json.Unmarshal(rr.Body.Bytes(), &out)
		require.NoError(t, err)

		assert.Equal(t, 2, out.TotalCount)
		assert.Equal(t, 1, len(out.Data))
		assert.Equal(t, testID1, out.Data[0].InstanceID)
		assert.True(t, out.PageInfo.HasNextPage)

		// given
		urlPath := fmt.Sprintf("/runtimes?cursor=%s", out.PageInfo.EndCursor)
		req, err = http.NewRequest(http.MethodGet, urlPath, nil)
		require.NoError(t, err)
		rr = httptest.NewRecorder()

		// when
		router.ServeHTTP(rr, req)

		// then
		require.Equal(t, http.StatusOK, rr.Code)

		err = json.Unmarshal(rr.Body.Bytes(), &out)
		require.NoError(t, err)
		logrus.Print(out.Data)
		assert.Equal(t, 2, out.TotalCount)
		assert.Equal(t, 1, len(out.Data))
		assert.Equal(t, testID2, out.Data[0].InstanceID)
		assert.False(t, out.PageInfo.HasNextPage)

	})

	t.Run("test validation should work", func(t *testing.T) {
		// given
		operations := memory.NewOperation()
		instances := memory.NewInstance(operations)

		runtimeHandler := NewHandler(instances, operations, 2, NewConverter())

		req, err := http.NewRequest("GET", "/runtimes?limit=a", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		runtimeHandler.AttachRoutes(router)

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		req, err = http.NewRequest("GET", "/runtimes?limit=1,2,3", nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		req, err = http.NewRequest("GET", "/runtimes?cursor=abc", nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

}
