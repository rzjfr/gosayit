package dict

import (
	"fmt"

	"github.com/hermanschaaf/enchant"
)

func Spell(word string) error {
	enchant, err := enchant.NewEnchant()
	if err != nil {
		return err
	}

	defer enchant.Free()

	// check whether a certain dictionary exists on the system
	has_en := enchant.DictExists("en_GB")

	// load the english dictionary:
	if has_en {
		enchant.LoadDict("en_GB")
		ok := enchant.Check(word)
		suggest := enchant.Suggest(word)
		if !ok {
			return fmt.Errorf("Cannot find the  word: %s\nDid you mean:%v", word, suggest)
		}
		return err
	} else {
		return fmt.Errorf("Problem with loading en_GB dictionary")
	}
}
