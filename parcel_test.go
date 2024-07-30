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
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, id, 0)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	newParsel, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, parcel.Address, newParsel.Address)
	assert.Equal(t, parcel.Client, newParsel.Client)
	assert.Equal(t, parcel.Status, newParsel.Status)
	assert.Equal(t, parcel.CreatedAt, newParsel.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)

	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// db, err := // настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")

	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, id, 0)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)

	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	newParsel, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, newAddress, newParsel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// db, err := // настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")

	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, id, 0)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := "new test status"

	err = store.SetStatus(id, newStatus)

	require.NoError(t, err)

	newParsel, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, newStatus, newParsel.Status)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

		require.NoError(t, err)

		assert.NotEqual(t, id, 0)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client

	require.NoError(t, err) // убедитесь в отсутствии ошибки

	assert.Equal(t, len(parcels), len(storedParcels)) // убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	// for _, parcel := range storedParcels {
	// 	// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
	// 	// убедитесь, что все посылки из storedParcels есть в parcelMap
	// 	// убедитесь, что значения полей полученных посылок заполнены верно
	// 	assert.Equal(t, parcel, parcelMap[client])
	// 	assert.Equal(t, parcel.Address, parcelMap[client].Address)
	// 	assert.Equal(t, parcel.Status, parcelMap[client].Status)
	// 	assert.Equal(t, parcel.CreatedAt, parcelMap[client].CreatedAt)
	// 	assert.Equal(t, parcel.Number, parcelMap[client].Number)
	// }

	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		assert.Equal(t, parcel.Number, parcelMap[parcel.Number].Number)
		assert.Equal(t, parcel.Address, parcelMap[parcel.Number].Address)
		assert.Equal(t, parcel.Status, parcelMap[parcel.Number].Status)
		assert.Equal(t, parcel.CreatedAt, parcelMap[parcel.Number].CreatedAt)
	}
}
