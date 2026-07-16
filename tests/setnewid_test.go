package unixid__test

import (
	. "github.com/tinywasm/unixid"
	"testing"

	. "github.com/tinywasm/fmt"
)

func TestSetNewID(t *testing.T) {
	// Creamos una sola instancia de UnixID para todos los subtests
	// para evitar la sobrecarga de crear múltiples instancias
	uid, err := NewUnixID()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("SetNewID con string", func(t *testing.T) {
		var id string
		uid.SetNewID(&id)

		if id == "" {
			t.Fatal("El ID generado no puede estar vacío")
		}

		// Validamos que tenga un formato correcto para servidor
		if Contains(id, ".") {
			t.Fatalf("En entorno servidor, el ID no debe contener punto: %s", id)
		}
	})

	t.Run("SetNewID con campo de struct", func(t *testing.T) {
		type User struct{ ID string }
		user := User{}
		uid.SetNewID(&user.ID)

		if user.ID == "" {
			t.Fatal("El ID generado no puede estar vacío")
		}
	})

	t.Run("Compatibilidad entre NewID y SetNewID", func(t *testing.T) {
		// Obtenemos ID con NewID
		idFromGet := uid.NewID()

		// Obtenemos ID con SetNewID
		var idFromSet string
		uid.SetNewID(&idFromSet)

		// Solo verificamos que ambos IDs tengan el mismo formato (longitud similar)
		lenGet := len(idFromGet)
		lenSet := len(idFromSet)
		if lenGet < lenSet-2 || lenGet > lenSet+2 { // Permitimos una pequeña variación
			t.Fatalf("Los IDs generados por NewID y SetNewID tienen formatos muy diferentes: %d vs %d", lenGet, lenSet)
		}
	})

	// Esta prueba solo funcionaría en compilación para WebAssembly
	// Se mantiene como referencia pero se omite en la ejecución
	t.Run("WebAssembly user number format (referencial)", func(t *testing.T) {
		t.Skip("Esta prueba está destinada para entornos WebAssembly")
	})
}

// Añadimos un benchmark para SetNewID para medir su rendimiento
func BenchmarkSetNewID(b *testing.B) {
	uid, _ := NewUnixID()
	var id string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid.SetNewID(&id)
	}
}
