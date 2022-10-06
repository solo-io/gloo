package dlp

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
)

var (
	SSN_REGEX_1                 = `(?:^|\D)([0-9]{9})(?:\D|$)`
	SSN_REGEX_2                 = `(?:^|\D)([0-9]{3}\-[0-9]{2}\-[0-9]{4})(?:\D|$)`
	SSN_REGEX_3                 = `(?:^|\D)([0-9]{3}\ [0-9]{2}\ [0-9]{4})(?:\D|$)`
	VISA_REGEX_1                = `(?:^|\D)(4[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	MASTERCARD_REGEX_1          = `(?:^|\D)(5[1-5][0-9]{2}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	DISCOVER_REGEX_1            = `(?:^|\D)(6011(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	AMEX_REGEX_1                = `(?:^|\D)((?:34|37)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{5})(?:\D|$)`
	JCB_REGEX_1                 = `(?:^|\D)(3[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	JCB_REGEX_2                 = `(?:^|\D)((?:2131|1800)[0-9]{11})(?:\D|$)`
	DINERS_CLUB_REGEX_1         = `(?:^|\D)(30[0-5][0-9](?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	DINERS_CLUB_REGEX_2         = `(?:^|\D)((?:36|38)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	CREDIT_CARD_TRACKER_REGEX_1 = `([1-9][0-9]{2}\-[0-9]{2}\-[0-9]{4}\^\d)`
	CREDIT_CARD_TRACKER_REGEX_2 = `(?:^|\D)(\%?[Bb]\d{13,19}\^[\-\/\.\w\s]{2,26}\^[0-9][0-9][01][0-9][0-9]{3})`
	CREDIT_CARD_TRACKER_REGEX_3 = `(?:^|\D)(\;\d{13,19}\=(?:\d{3}|)(?:\d{4}|\=))`

	ssnTransform = &transformation_ee.Action{
		Name: "ssn",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    SSN_REGEX_1,
				Subgroup: 1,
			},
			{
				Regex:    SSN_REGEX_2,
				Subgroup: 1,
			},
			{
				Regex:    SSN_REGEX_3,
				Subgroup: 1,
			},
		},
	}

	visaTransform = &transformation_ee.Action{
		Name: "visa",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    VISA_REGEX_1,
				Subgroup: 1,
			},
		},
	}

	mastercardTransform = &transformation_ee.Action{
		Name: "master_card",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    MASTERCARD_REGEX_1,
				Subgroup: 1,
			},
		},
	}

	discoverTransform = &transformation_ee.Action{
		Name: "discover",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    DISCOVER_REGEX_1,
				Subgroup: 1,
			},
		},
	}

	amexTransform = &transformation_ee.Action{
		Name: "amex",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    AMEX_REGEX_1,
				Subgroup: 1,
			},
		},
	}

	jcbTransform = &transformation_ee.Action{
		Name: "jcb",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    JCB_REGEX_1,
				Subgroup: 1,
			},
			{
				Regex:    JCB_REGEX_2,
				Subgroup: 1,
			},
		},
	}

	dinersClubTransform = &transformation_ee.Action{
		Name: "diners_club",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    DINERS_CLUB_REGEX_1,
				Subgroup: 1,
			},
			{
				Regex:    DINERS_CLUB_REGEX_2,
				Subgroup: 1,
			},
		},
	}

	creditCardTrackersTransform = &transformation_ee.Action{
		Name: "credit_card_trackers",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex:    CREDIT_CARD_TRACKER_REGEX_1,
				Subgroup: 1,
			},
			{
				Regex:    CREDIT_CARD_TRACKER_REGEX_2,
				Subgroup: 1,
			},
			{
				Regex:    CREDIT_CARD_TRACKER_REGEX_3,
				Subgroup: 1,
			},
		},
	}

	allCreditCardTransform = &transformation_ee.Action{
		Name: "all_credit_cards",
		RegexActions: []*transformation_ee.RegexAction{
			{
				Regex: "(" + strings.Join([]string{
					VISA_REGEX_1,
					MASTERCARD_REGEX_1,
					DISCOVER_REGEX_1,
					AMEX_REGEX_1,
					JCB_REGEX_1,
					JCB_REGEX_2,
					DINERS_CLUB_REGEX_1,
					DINERS_CLUB_REGEX_2,
					CREDIT_CARD_TRACKER_REGEX_1,
					CREDIT_CARD_TRACKER_REGEX_2,
					CREDIT_CARD_TRACKER_REGEX_3,
				}, "|") + ")",
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
		dlp.Action_ALL_CREDIT_CARDS_COMBINED: {allCreditCardTransform},
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
