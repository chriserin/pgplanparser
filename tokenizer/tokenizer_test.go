package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizerSimple(t *testing.T) {

	plan := []rune("{SOMETHING :keyA.something valueB}")
	tokens := Tokenize(plan)

	assert.Equal(t, ItemStart, tokens[0].Token)
	assert.Equal(t, ItemId, tokens[1].Token)
	assert.Equal(t, ItemKey, tokens[2].Token)
	assert.Equal(t, ItemValue, tokens[3].Token)
	assert.Equal(t, ItemEnd, tokens[4].Token)
}

func TestTokenizerList(t *testing.T) {

	plan := []rune("{SOMETHING :keyA.something (a 1)}")
	tokens := Tokenize(plan)

	assert.Equal(t, ItemStart, tokens[0].Token)
	assert.Equal(t, 1, tokens[0].Depth)
	assert.Equal(t, ItemId, tokens[1].Token)
	assert.Equal(t, 1, tokens[1].Depth)
	assert.Equal(t, ItemKey, tokens[2].Token)
	assert.Equal(t, ListStart, tokens[3].Token)
	assert.Equal(t, ListValue, tokens[4].Token)
	assert.Equal(t, ListValue, tokens[5].Token)
	assert.Equal(t, 2, tokens[5].Depth)
	assert.Equal(t, ListEnd, tokens[6].Token)
	assert.Equal(t, 2, tokens[6].Depth)
	assert.Equal(t, ItemEnd, tokens[7].Token)
	assert.Equal(t, 1, tokens[7].Depth)
}

func TestTokenizerFull(t *testing.T) {

	plan := []rune(`
  {PLANNEDSTMT :commandType 1 :queryId 16893614937036654096 :hasReturning false
        :hasModifyingCTE false :canSetTag true :transientPlan false :dependsOnRole
        false :parallelModeNeeded false :jitFlags 0 :planTree {SEQSCAN
        :scan.plan.startup_cost 0 :scan.plan.total_cost 15455.779999999999
        :scan.plan.plan_rows 683178 :scan.plan.plan_width 4 :scan.plan.parallel_aware
        false :scan.plan.parallel_safe true :scan.plan.async_capable false
        :scan.plan.plan_node_id 0 :scan.plan.targetlist ({TARGETENTRY :expr {VAR
        :varno 1 :varattno 1 :vartype 23 :vartypmod -1 :varcollid 0 :varnullingrels
        (b) :varlevelsup 0 :varnosyn 1 :varattnosyn 1 :location 7} :resno 1 :resname
        flight_id :ressortgroupref 0 :resorigtbl 16424 :resorigcol 1 :resjunk false})
        :scan.plan.qual <> :scan.plan.lefttree <> :scan.plan.righttree <>
        :scan.plan.initPlan <> :scan.plan.extParam (b) :scan.plan.allParam (b)
        :scan.scanrelid 1} :rtable ({RANGETBLENTRY :alias <> :eref {ALIAS :aliasname
        flight :colnames ("flight_id" "flight_no" "scheduled_departure"
        "scheduled_arrival" "departure_airport" "arrival_airport" "status"
        "aircraft_code" "actual_departure" "actual_arrival" "update_ts")} :rtekind 0
        :relid 16424 :relkind r :rellockmode 1 :tablesample <> :perminfoindex 1
        :lateral false :inh false :inFromCl true :securityQuals <>}) :permInfos
        ({RTEPERMISSIONINFO :relid 16424 :inh true :requiredPerms 2 :checkAsUser 0
        :selectedCols (b 8) :insertedCols (b) :updatedCols (b)}) :resultRelations <>
        :appendRelations <> :subplans <> :rewindPlanIDs (b) :rowMarks <> :relationOids
        (o 16424) :invalItems <> :paramExecTypes <> :utilityStmt <> :stmt_location 0
        :stmt_len 28}
  `)
	tokens := Tokenize(plan)

	assert.Equal(t, tokens[0].Token, ItemStart)
	assert.Equal(t, 205, len(tokens))
}
