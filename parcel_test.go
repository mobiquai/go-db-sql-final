package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err) // завершить тест, если функция вернула ошибку
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку
	require.NotEmpty(t, id) // завершить тест, если пустое

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	testGet, err := store.Get(id)
	require.NoError(t, err)          // завершить тест, если функция вернула ошибку
	assert.Equal(t, parcel, testGet) // если значения всех полей полученного объекта НЕ совпадают со значениями полей объекта в переменной parcel, то ошибка

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку

	testGet, err = store.Get(id)
	require.NoError(t, err)  // завершить тест, если функция вернула ошибку
	assert.Empty(t, testGet) // если структура НЕ пустая (посылку больше нельзя получить из БД), то ошибка

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err) // завершить тест, если функция вернула ошибку
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку
	require.NotEmpty(t, id) // завершить тест, если пустое

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	testGet, err := store.Get(id)
	require.NoError(t, err)                      // завершить тест, если функция вернула ошибку
	assert.Equal(t, newAddress, testGet.Address) // если значение newAddress НЕ совпадают со значением Address в полученной записи БД, то ошибка

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err) // завершить тест, если функция вернула ошибку
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку
	require.NotEmpty(t, id) // завершить тест, если пустое

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err) // завершить тест, если функция вернула ошибку

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	testGet, err := store.Get(id)
	require.NoError(t, err)                           // завершить тест, если функция вернула ошибку
	assert.Equal(t, ParcelStatusSent, testGet.Status) // если значение ParcelStatusSent НЕ совпадают со значением Status в полученной записи БД, то ошибка

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err) // завершить тест, если функция вернула ошибку
	}
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int64]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		require.NoError(t, err) // завершить тест, если функция вернула ошибку
		require.NotEmpty(t, id) // завершить тест, если пустое

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)                           // завершить тест, если функция вернула ошибку
	assert.Equal(t, len(parcels), len(storedParcels)) // если значение len(parcels) НЕ совпадают со значением len(storedParcels), то ошибка

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно

		_, isParcel := parcelMap[parcel.Number] // проверка наличия элемента в мапе по ключу
		require.True(t, isParcel)               // Если текущей посылки из storedParcels нет в parcelMap, то ошибка

		assert.Equal(t, parcelMap[parcel.Number], parcel) // если значения всех полей полученного объекта НЕ совпадают с найденной посылкой из мапы, то ошибка
	}
}
