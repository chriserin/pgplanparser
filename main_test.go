package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainAgain(t *testing.T) {
	os.Args = []string{"exename", "{PLANNEDSTMT }"}
	main()
}

func TestParsePlanSimple(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SEQSCAN lefttree {LOOPA} junk ({NOTHING x 1}) righttree {LOOPB}}}"

	value := processPlan(planDetail)

	assert.Equal(t, "SEQSCAN", value.Plantree.Nodetype)
	assert.Equal(t, "LOOPA", value.Plantree.Lefttree.Nodetype)
	assert.Equal(t, "LOOPB", value.Plantree.Righttree.Nodetype)
}

func TestParseRtables(t *testing.T) {

	planDetail := "{PLANNEDSTMT :rtables ({RTABLEENTRY relid 16424 junk ({NOTHING x 1})})}"

	value := processPlan(planDetail)

	assert.Equal(t, 16424, value.Rtables[0].Relid)
}

func TestParseAssignTableEntryId(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SEQSCAN relid 1} :rtables ({RTABLEENTRY relid 16424 })}"

	value := processPlan(planDetail)

	assert.Equal(t, 1, value.Plantree.Relid)
}

func TestParseAssignTableName(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SEQSCAN relid 1} :rtables ({RTABLEENTRY relid 16424 })}"

	value := processPlan(planDetail)
	populateTableNames(&value, []postgresTable{{16424, "flight"}})

	assert.Equal(t, "flight", value.Plantree.Tablename)
}

func TestParseSetopExcept(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SETOP cmd 2}}"

	value := processPlan(planDetail)

	assert.Equal(t, 2, value.Plantree.Cmd)
}

func TestParseSetopIntersect(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SETOP cmd 0}}"

	value := processPlan(planDetail)

	assert.Equal(t, 0, value.Plantree.Cmd)
}

func TestParseSetopHash(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SETOP strategy 1}}"

	value := processPlan(planDetail)

	assert.Equal(t, 1, value.Plantree.Strategy)
}

func TestParseSetopHashIntersect(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SETOP cmd 0 strategy 1}}"

	value := processPlan(planDetail)

	assert.Equal(t, 0, value.Plantree.Cmd)
	assert.Equal(t, 1, value.Plantree.Strategy)
}

func TestParseAssignTableNameOnNestedNode(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SEQSCAN lefttree {SEQSCAN relid 1}} :rtables ({RTABLEENTRY relid 16424 })}"

	value := processPlan(planDetail)
	populateTableNames(&value, []postgresTable{{16424, "flight"}})

	assert.Equal(t, "flight", value.Plantree.Lefttree.Tablename)
}

func TestParseNullLeftTree(t *testing.T) {

	planDetail := "{PLANNEDSTMT :planTree {SEQSCAN lefttree <>}}"

	value := processPlan(planDetail)

	assert.Equal(t, "SEQSCAN", value.Plantree.Nodetype)
	assert.Equal(t, true, value.Plantree.Lefttree == nil)
}

func TestParsePlanSelect1(t *testing.T) {

	// select 1;
	planDetail := `{PLANNEDSTMT :commandType 1 :queryId 6865378226349601843 :hasReturning false
        :hasModifyingCTE false :canSetTag true :transientPlan false :dependsOnRole
        false :parallelModeNeeded false :jitFlags 0 :planTree {RESULT
        :plan.startup_cost 0 :plan.total_cost 0.01 :plan.plan_rows 1 :plan.plan_width
        4 :plan.parallel_aware false :plan.parallel_safe false :plan.async_capable
        false :plan.plan_node_id 0 :plan.targetlist ({TARGETENTRY :expr {CONST
        :consttype 23 :consttypmod -1 :constcollid 0 :constlen 4 :constbyval true
        :constisnull false :location 7 :constvalue 4 [ 1 0 0 0 0 0 0 0 ]} :resno 1
        :resname ?column? :ressortgroupref 0 :resorigtbl 0 :resorigcol 0 :resjunk
        false}) :plan.qual <> :plan.lefttree <> :plan.righttree <> :plan.initPlan <>
        :plan.extParam (b) :plan.allParam (b) :resconstantqual <>} :rtable
        ({RANGETBLENTRY :alias <> :eref {ALIAS :aliasname *RESULT* :colnames <>}
        :rtekind 8 :lateral false :inh false :inFromCl false :securityQuals <>})
        :permInfos <> :resultRelations <> :appendRelations <> :subplans <>
        :rewindPlanIDs (b) :rowMarks <> :relationOids <> :invalItems <>
        :paramExecTypes <> :utilityStmt <> :stmt_location 0 :stmt_len 8}`

	value := processPlan(planDetail)
	assert.NotNil(t, value)
}

func TestParsePlanSelectFromTable(t *testing.T) {

	// Select * from flight;
	planDetail := `{PLANNEDSTMT :commandType 1 :queryId 16893614937036654096 :hasReturning false
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
        :stmt_len 28}`

	value := processPlan(planDetail)
	populateTableNames(&value, []postgresTable{{16424, "flight"}})

	assert.Equal(t, "flight", value.Plantree.Tablename)
}
