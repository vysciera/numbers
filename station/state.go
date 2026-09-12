package station

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const stateVersion = 1

type State struct {
	Version	uint8 `json:"version"`

	StationID	uint16	`json:"station_id"`
	Epoch		string 	`json:"epoch"`

	Sequence	uint32 `json:"sequence"`
	Transmission uint32 `json:"transmission"`

	CreatedAt	time.Time `json:"created_at"`
	LastStartedAt	time.Time `json:"last_started_at"`
	Starts	uint64 `json:"starts"`
}

func newState(stationID uint16) (State, error) {
	if stationID == 0 {
		return State{}, errors.New("station ID 0 is reserved")
	}

	epoch, err := newEpoch()
	if err != nil {
		return State{}, err
	}

	now := time.Now().UTC()

	return State{
		Version:	stateVersion,
		StationID:	stationID,
		Epoch:		epoch,
		CreatedAt:	now,
		LastStartedAt:	now,
		Starts:	1,
	}, nil
}

func newEpoch() (string, error) {
	var raw [16]byte

	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate station epoch: %w", err)
	}

	return hex.EncodeToString(raw[:]), nil
}

func loadState(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return State{}, err
	}

	var state State

	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("decode station state: %w", err)
	}

	if state.Version != stateVersion {
		return State{}, fmt.Errorf("unsupported state version: %w", state.Version)
	}

	if state.StationID == 0 {
		return State{}, errors.New("invalid persisted station ID")
	}

	if state.Epoch == "" {
		return State{}, errors.New("persisted state has no epoch")
	}

	return state, nil
}

func saveState(path string, state State) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	data, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return fmt.Errorf("encode station state: %w", err)
	}

	data = append(data, '\n')

	tmp, err := os.CreateTemp(dir, ".numbers-state-*")
	if err != nil {
		return fmt.Errorf("create temporary state: %w", err)
	}

	tmpName := tmp.Name()

	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	if err := tmp.Chmod(0644); err != nil {
		return fmt.Errorf("chmod temporary state: %w", err)
	}

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("write temporary state: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temporary state: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary state: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace station state: %w", err)
	}

	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}

	return nil
}
