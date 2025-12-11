package unixid

import (
	"testing"
	"time"

	"github.com/tinywasm/unixid"
)

// TestGetNewIDWithCorrectFormatting prueba el flujo completo de generación de IDs
func TestGetNewIDWithCorrectFormatting(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	// Simular generación de múltiples IDs
	var ids []string

	for i := 0; i < 3; i++ {
		id := uid.GetNewID() // Devuelve nanosegundos como string
		ids = append(ids, id)

		t.Logf("Mensaje %d - ID: %s", i+1, id)

		// Pausa de 1 segundo para garantizar diferencias de tiempo visibles
		time.Sleep(1 * time.Second)
	}

	// Verificar orden cronológico de IDs
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("Los IDs NO están en orden cronológico: %s <= %s",
				ids[i], ids[i-1])
		}
	}

	// t.Log("✅ IDs están en orden cronológico correcto")
}
