package tv

import (
	"fmt"
	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"log"
	"strings"
)

type shortcut struct {
	r      rune
	key    tcell.Key
	mod    tcell.ModMask
	action view.Action
}

type shortcutModActions map[tcell.ModMask]view.Action

type shortcutKeyMap map[tcell.Key]shortcutModActions

type shortcutRuneMap map[rune]shortcutModActions

type shortcutMap struct {
	keyMap  shortcutKeyMap
	runeMap shortcutRuneMap
}

func (s *shortcut) keyName() string {
	result := ""
	m := []string{}
	if s.mod & tcell.ModShift != 0 {
		m = append(m, "Shift")
	}
	if s.mod & tcell.ModAlt != 0 {
		m = append(m, "Alt")
	}
	if s.mod & tcell.ModMeta != 0 {
		m = append(m, "Meta")
	}
	if s.mod & tcell.ModCtrl != 0 {
		m = append(m, "Ctrl")
	}

	key := s.key
	if s.r != 0 {
		key = tcell.KeyRune
	}
	ok := false
	if result, ok = tcell.KeyNames[key]; !ok {
		if key == tcell.KeyRune {
			if s.r == ' ' {
				result = "space"
			} else {
				result = string(s.r)
			}
		} else {
			result = fmt.Sprintf("Key[%d,%d]", key, int(s.r))
		}
	}
	if len(m) != 0 {
		if s.mod & tcell.ModCtrl != 0 && strings.HasPrefix(result, "Ctrl-") {
			result = result[5:]
		}
		return fmt.Sprintf("%s+%s", strings.Join(m, "+"), result)
	}
	return result
}


func generateActionShortcutNames(shortcuts []shortcut) map[view.Action][]string {
	result := make(map[view.Action][]string)
	for _, s := range shortcuts {
		if _, exists := result[s.action]; !exists {
			result[s.action] = nil
		}
		result[s.action] = append(result[s.action], s.keyName())
	}
	return result
}




func newShortcutMap(shortcuts []shortcut) *shortcutMap {
	result := &shortcutMap{
		keyMap:  make(shortcutKeyMap),
		runeMap: make(shortcutRuneMap),
	}

	getRuneMap := func(r rune) shortcutModActions {
		if _, exists := result.runeMap[r]; !exists {
			result.runeMap[r] = make(shortcutModActions)
		}
		return result.runeMap[r]
	}

	getKeyMap := func(k tcell.Key) shortcutModActions {
		if _, exists := result.keyMap[k]; !exists {
			result.keyMap[k] = make(shortcutModActions)
		}
		return result.keyMap[k]
	}

	for _, s := range shortcuts {
		var modMap shortcutModActions
		if s.r != 0 { // Rune
			if s.mod & tcell.ModShift != 0 {
				log.Panicf("Rune (%c) shortcut specified with shift modifier.", s.r)
			}
			modMap = getRuneMap(s.r)
		} else { // Key
			if keyName, exists := tcell.KeyNames[s.key]; exists {
				if strings.HasPrefix(keyName, "Ctrl-") {
					s.mod |= tcell.ModCtrl
				}
			}
			modMap = getKeyMap(s.key)
		}
		if _, exists := modMap[s.mod]; exists {
			log.Panicf("Duplicated shortcut (rune %d, key %d, mod %d).", s.r, s.key, s.mod)
		}
		modMap[s.mod] = s.action
	}

	return result
}


func (sm *shortcutMap) mapKeys(ev *tcell.EventKey) view.Action {
	result := view.ActionUnknown
	key := ev.Key()
	mod := ev.Modifiers()
	if key == tcell.KeyRune {
		if m, exists := sm.runeMap[ev.Rune()]; exists {
			if v, exists := m[mod]; exists {
				result = v
			}
		}
	} else {
		if m, exists := sm.keyMap[key]; exists {
			if v, exists := m[mod]; exists {
				result = v
			}
		}
	}
	return result
}



