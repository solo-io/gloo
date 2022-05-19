package dlp

import (
	"github.com/golang/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
)

var (
	ssnTransform = &transformation_ee.Action{
		Name: "ssn",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)([0-9]{9})(?:\D|$)`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)([0-9]{3}\-[0-9]{2}\-[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)([0-9]{3}\ [0-9]{2}\ [0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	visaTransform = &transformation_ee.Action{
		Name: "visa",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)(4[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	mastercardTransform = &transformation_ee.Action{
		Name: "master_card",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)(5[1-5][0-9]{2}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	discoverTransform = &transformation_ee.Action{
		Name: "discover",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)(6011(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	amexTransform = &transformation_ee.Action{
		Name: "amex",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)((?:34|37)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{5})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	jcbTransform = &transformation_ee.Action{
		Name: "jcb",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)(3[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)((?:2131|1800)[0-9]{11})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	dinersClubTransform = &transformation_ee.Action{
		Name: "diners_club",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `(?:^|\D)(30[0-5][0-9](?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)((?:36|38)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`,
				Subgroup: 1,
			},
		},
	}

	creditCardTrackersTransform = &transformation_ee.Action{
		Name: "credit_card_trackers",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    `([1-9][0-9]{2}\-[0-9]{2}\-[0-9]{4}\^\d)`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)(\%?[Bb]\d{13,19}\^[\-\/\.\w\s]{2,26}\^[0-9][0-9][01][0-9][0-9]{3})`,
				Subgroup: 1,
			},
			{
				Regex:    `(?:^|\D)(\;\d{13,19}\=(?:\d{3}|)(?:\d{4}|\=))`,
				Subgroup: 1,
			},
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
