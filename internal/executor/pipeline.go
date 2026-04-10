package executor

import 

// RunTestcases calls the judge layer and maps the failed key
// back onto the aggregated result for the worker to persist.
func (p *Pipeline) RunTestcases(
	tcMap *judge.TestCaseMap,
	ws *Workspace,
	lang models.Language,
) ([]judge.RunResult, string, error) {

	results, failedKey, err := judge.RunAllTestcases(tcMap, ws, lang)
	if err != nil {
		// failedKey is already populated — caller logs it and
		// stores it on the submission for the client to read back
		return results, failedKey, err
	}

	return results, "", nil
}