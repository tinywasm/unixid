package unixid

import (
	"sync"
	"testing"
	"time"

	"github.com/tinywasm/unixid"
)

func Test_GetNewID(t *testing.T) {
	idRequired := 10000
	wg := sync.WaitGroup{}
	wg.Add(idRequired)

	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	idObtained := make(map[string]int)
	var esperar sync.Mutex

	for i := 0; i < idRequired; i++ {
		go func() {
			defer wg.Done()

			id := uid.GetNewID()

			esperar.Lock()
			if cantId, exist := idObtained[id]; exist {
				idObtained[id] = cantId + 1
			} else {
				idObtained[id] = 1
			}
			esperar.Unlock()

		}()
	}
	wg.Wait()

	// fmt.Printf("total id requeridos: %v ob: %v\n", idRequired, len(idObtained))
	if idRequired != len(idObtained) {
		t.Fatalf("se esperaban: %d ids pero se obtuvieron: %d. Detalle: %v", idRequired, len(idObtained), idObtained)
	}

}

func BenchmarkGetNewID(b *testing.B) {
	uid, _ := unixid.NewUnixID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid.GetNewID()
	}
}

// Prueba adicional para verificar que no haya duplicados al generar muchos IDs
func TestNoDuplicateIDs(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Generar una cantidad moderada de IDs y verificar que no haya duplicados
	numIDs := 1000
	ids := make(map[string]bool)

	for i := 0; i < numIDs; i++ {
		id := uid.GetNewID()
		if _, exists := ids[id]; exists {
			t.Fatalf("ID duplicado encontrado: %s", id)
		}
		ids[id] = true
	}
}

// Prueba para verificar que se generen IDs secuenciales cuando hay colisiones de timestamp
func TestSequentialIDs(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Generar varios IDs rápidamente, algunos tendrán el mismo timestamp base
	// pero deberían tener números secuenciales añadidos
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		ids[i] = uid.GetNewID()
	}

	// Verificar que tengamos al menos algunos IDs diferentes
	uniqueIDs := make(map[string]bool)
	for _, id := range ids {
		uniqueIDs[id] = true
	}

	if len(uniqueIDs) < len(ids) {
		t.Fatalf("Se esperaban %d IDs únicos, pero se obtuvieron %d", len(ids), len(uniqueIDs))
	}
}

// TestExternalMutexNoDeadlock verifica que cuando se proporciona un mutex externo,
// la llamada a GetNewID no se bloquee cuando ya hay un lock adquirido con el mismo mutex.
// Esta prueba verifica el comportamiento actualizado donde se usa un no-op mutex internamente
// cuando se proporciona un mutex externo para prevenir deadlocks.
func TestExternalMutexNoDeadlock(t *testing.T) {
	// Creamos un mutex externo que simula ser compartido con otra biblioteca
	externalMutex := &sync.Mutex{}

	// Creamos una instancia de UnixID pasando el mutex externo
	uid, err := unixid.NewUnixID(externalMutex)
	if err != nil {
		t.Fatalf("Error creando UnixID con mutex externo: %v", err)
		return
	}

	// Simulamos un escenario donde otra biblioteca bloquea el mutex
	// y luego nuestro código lo utiliza
	externalMutex.Lock()
	defer externalMutex.Unlock()

	// Definimos un canal para detectar si hay deadlock
	done := make(chan bool)
	go func() {
		// Esto NO debería bloquearse con la nueva lógica, ya que
		// internamente estamos usando un defaultNoOpMutex
		id := uid.GetNewID()
		if id == "" {
			t.Error("Se generó un ID vacío")
		}
		done <- true
	}()

	// Esperamos brevemente para ver si se completa la generación del ID
	select {
	case <-done:
		// Este es el comportamiento esperado: GetNewID no se bloquea
		// porque internamente estamos usando un defaultNoOpMutex
	case <-time.After(time.Millisecond * 500):
		t.Fatal("GetNewID se bloqueó a pesar de usar un no-op mutex internamente")
	}

	// Verificación adicional: generar varios IDs sin problemas mientras el mutex está bloqueado
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id := uid.GetNewID()
		if id == "" {
			t.Fatalf("Se generó un ID vacío en la iteración %d", i)
		}

		if _, exists := ids[id]; exists {
			t.Fatalf("ID duplicado encontrado: %s", id)
		}
		ids[id] = true
	}
}
