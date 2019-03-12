package services

import (
	"github.com/lca1/drynx/lib"
	"github.com/lca1/drynx/protocols"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

// Query Handlers
//______________________________________________________________________________________________________________________

// HandleSurveyQueryToDP handles the reception of a query at a DP
func (s *ServiceDrynx) HandleSurveyQueryToDP(recq *libdrynx.SurveyQueryToDP) (network.Message, error) {

	recq.SQ.Query.IVSigs.InputValidationSigs = recreateRangeSignatures(recq.SQ.Query.IVSigs)
	// only generate ProofCollection protocol instances if proofs is enabled
	var mapPIs map[string]onet.ProtocolInstance
	if recq.SQ.Query.Proofs != 0 {
		mapPIs = s.generateRangePI(recq)
	}

	s.Survey.Put(recq.SQ.SurveyID, Survey{
		SurveyQuery: recq.SQ,
		MapPIs:      mapPIs,
	})

	// signal the root that the data provider has received the query
	err := s.SendRaw(recq.Root, &DPqueryReceived{recq.SQ.SurveyID})
	if err != nil {
		log.Error("[SERVICE] <Drynx> Server, broadcasting [DPdataFinished] error ", err)
	}

	return nil, nil
}

// Support Functions
//______________________________________________________________________________________________________________________

func (s *ServiceDrynx) generateRangePI(query *libdrynx.SurveyQueryToDP) map[string]onet.ProtocolInstance {
	mapPIs := make(map[string]onet.ProtocolInstance)
	for _, dp := range *query.SQ.ServerToDP[query.Root.String()] {
		if dp.String() == s.ServerIdentity().String() {
			tree := generateProofCollectionRoster(&dp, query.SQ.Query.RosterVNs).GenerateStar()

			pi := s.CreateProofCollectionPIs(tree, query.SQ.SurveyID, protocols.ProofCollectionProtocolName)
			mapPIs["range/"+s.ServerIdentity().String()] = pi
			break
		}

	}

	return mapPIs

}
