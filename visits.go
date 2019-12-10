package gostp

import (
	"encoding/gob"
	"errors"
	"net/http"
	"os"

	"github.com/axiomhq/hyperloglog"
	"github.com/tomasen/realip"
)

var CounterMap map[string]*hyperloglog.Sketch

// ErrCount - error returned when you try to get count but didn't register middleware
var ErrCount = errors.New("count not found or error in HyperLogLog")

// Visits - get visits for given URL
func Visits(r *http.Request) (uint64, error) {
	if CounterMap == nil {
		// no, you didn't ...
		panic("you need to register Visigo Counter first!")
	}
	if hll, found := CounterMap[r.URL.String()]; found {
		return hll.Estimate(), nil
	}
	return 0, ErrCount
}

// TotalVisits gets total visits to all sites
func TotalVisits() (uint64, error) {
	hll := hyperloglog.New()
	for _, s := range CounterMap {
		if err := hll.Merge(s); err != nil {
			return 0, err
		}
	}
	return hll.Estimate(), nil
}

// Counter - registers middleware for visits counting
func Counter(next http.Handler) http.Handler {
	CounterMap = make(map[string]*hyperloglog.Sketch)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hll, found := CounterMap[r.URL.String()]; found {
			hll.Insert([]byte(realip.RealIP(r)))
		} else {
			l := hyperloglog.New()
			l.Insert([]byte(realip.RealIP(r)))
			CounterMap[r.URL.String()] = l
		}
		next.ServeHTTP(w, r)
	})
}

func SaveVisits() {
	//////////
	// First lets encode some data
	//////////

	// Create a file for IO
	var encodeFile *os.File
	if FileNotExist("visits.gob") == nil {
		encodeFile, _ = os.Create("visits.gob")
	} else {
		encodeFile, _ = os.Open("visits.gob")
	}

	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(CounterMap); err != nil {
		panic(err)
	}
	encodeFile.Close()
}
