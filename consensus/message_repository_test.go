package consensus

import (
	"testing"

	"sync"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMsgRepositoryImpl_Save(t *testing.T) {
	// Given
	prepareMsgRepository := NewPrepareMsgRepository()

	consensusId := NewConsensusId("1")

	prepareMsg := PrepareMsg{
		ConsensusId: consensusId,
		SenderId:    "1",
	}

	// When
	prepareMsgRepository.Save(prepareMsg)

	// Then
	assert.Equal(t, 1, len(prepareMsgRepository.FindPrepareMsgsByConsensusID(consensusId)))
}

func TestPrepareMsgRepositoryImpl_Remove(t *testing.T) {
	// Given
	prepareMsgRepository := NewPrepareMsgRepository()
	prepareMsgRepositoryImpl := PrepareMsgRepositoryImpl{
		PreparePool: make(map[ConsensusId][]PrepareMsg, 0),
		lock:        &sync.RWMutex{},
	}
	prepareMsgRepository = &prepareMsgRepositoryImpl

	consensusId := NewConsensusId("1")

	prepareMsg := PrepareMsg{
		ConsensusId: consensusId,
		SenderId:    "1",
	}

	// When
	prepareMsgRepository.Save(prepareMsg)
	prepareMsgRepository.Remove(consensusId)

	// Then
	assert.Equal(t, 0, len(prepareMsgRepositoryImpl.PreparePool))
}

func TestCommitMsgRepositoryImpl_Save(t *testing.T) {
	// Given
	prepareMsgRepository := NewPrepareMsgRepository()

	consensusId := NewConsensusId("1")

	prepareMsg := PrepareMsg{
		ConsensusId: consensusId,
		SenderId:    "1",
	}

	// When
	prepareMsgRepository.Save(prepareMsg)

	// Then
	assert.Equal(t, 1, len(prepareMsgRepository.FindPrepareMsgsByConsensusID(consensusId)))
}

func TestCommitMsgRepositoryImpl_Remove(t *testing.T) {
	// Given
	commitMsgRepository := NewCommitMsgRepository()
	commitMsgRepositoryImpl := CommitMsgRepositoryImpl{
		CommitPool: make(map[ConsensusId][]CommitMsg, 0),
		lock:       &sync.RWMutex{},
	}
	commitMsgRepository = &commitMsgRepositoryImpl

	consensusId := NewConsensusId("1")

	commitMsg := CommitMsg{
		ConsensusId: consensusId,
		SenderId:    "1",
	}

	// When
	commitMsgRepository.Save(commitMsg)
	commitMsgRepository.Remove(consensusId)

	// Then
	assert.Equal(t, 0, len(commitMsgRepositoryImpl.CommitPool))
}