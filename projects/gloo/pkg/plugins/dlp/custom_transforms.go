package dlp

import (
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/dlp"
)

var (
	ssnTransform = &transformation_ee.Action{
		Name: "ssn",
		Regex: []string{
			"(?!\\D)[0-9]{9}(?=\\D|$)",
			"(?!\\D)[0-9]{3}\\-[0-9]{2}\\-[0-9]{4}(?=\\D|$)",
			"(?!\\D)[0-9]{3}\\ [0-9]{2}\\ [0-9]{4}(?=\\D|$)",
		},
	}

	visaTransform = &transformation_ee.Action{
		Name: "visa",
		Regex: []string{
			"(?!\\D)4[0-9]{3}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
		},
	}

	mastercardTransform = &transformation_ee.Action{
		Name: "master_card",
		Regex: []string{
			"(?!\\D)5[1-5][0-9]{2}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
		},
	}

	discoverTransform = &transformation_ee.Action{
		Name: "discover",
		Regex: []string{
			"(?!\\D)6011(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
		},
	}

	amexTransform = &transformation_ee.Action{
		Name: "amex",
		Regex: []string{
			"(?!\\D)(34|37)[0-9]{2}(\\ |\\-|)[0-9]{6}(\\ |\\-|)[0-9]{5}(?=\\D|$)",
		},
	}

	jcbTransform = &transformation_ee.Action{
		Name: "jcb",
		Regex: []string{
			"(?!\\D)3[0-9]{3}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
			"(?!\\D)(2131|1800)[0-9]{11}(?=\\D|$)",
		},
	}

	dinersClubTransform = &transformation_ee.Action{
		Name: "diners_club",
		Regex: []string{
			"(?!\\D)30[0-5][0-9](\\ |\\-|)[0-9]{6}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
			"(?!\\D)(36|38)[0-9]{2}(\\ |\\-|)[0-9]{6}(\\ |\\-|)[0-9]{4}(?=\\D|$)",
		},
	}

	creditCardTrackersTransform = &transformation_ee.Action{
		Name: "credit_card_trackers",
		Regex: []string{
			"[1-9][0-9]{2}\\-[0-9]{2}\\-[0-9]{4}\\^\\d",
			"(?!\\D)\\%?[Bb]\\d{13,19}\\^[\\-\\/\\.\\w\\s]{2,26}\\^[0-9][0-9][01][0-9][0-9]{3}",
			"(?!\\D)\\;\\d{13,19}\\=(\\d{3}|)(\\d{4}|\\=)",
		},
	}

	transformMap = map[dlp.Action_ActionType][]*transformation_ee.Action{
		dlp.Action_SSN:                  {ssnTransform},
		dlp.Action_MASTERCARD:           {mastercardTransform},
		dlp.Action_VISA:                 {visaTransform},
		dlp.Action_AMEX:                 {amexTransform},
		dlp.Action_DISCOVER:             {discoverTransform},
		dlp.Action_JCB:                  {jcbTransform},
		dlp.Action_DINERS_CLUB:          {dinersClubTransform},
		dlp.Action_CREDIT_CARD_TRACKERS: {creditCardTrackersTransform},
		dlp.Action_ALL_CREDIT_CARDS: {
			mastercardTransform,
			visaTransform,
			amexTransform,
			discoverTransform,
			jcbTransform,
			dinersClubTransform,
			creditCardTrackersTransform,
		},
	}
)

func GetTransformsFromMap(actionType dlp.Action_ActionType) []*transformation_ee.Action {
	var result []*transformation_ee.Action
	transformers := transformMap[actionType]
	for _, v := range transformers {
		transformerMsg := proto.Clone(v).(*transformation_ee.Action)
		result = append(result, transformerMsg)
	}
	return result
}
