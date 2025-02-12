package gui

import (
	"github.com/jesseduffield/gocui"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{} // FIXME: find out how to get `gocui.Key | rune`
	Modifier    gocui.Modifier
	Description string
	Alternative string
}

// GetDisplayStrings returns the display string of a file
func (b *Binding) GetDisplayStrings(isFocused bool) []string {
	return []string{b.GetKey(), b.Description}
}

// GetKey is a function.
func (b *Binding) GetKey() string {
	key := 0

	switch b.Key.(type) {
	case rune:
		key = int(b.Key.(rune))
	case gocui.Key:
		if b.Key.(gocui.Key) == gocui.KeyCtrlJ {
			return "ctrl+j"
		}
		if b.Key.(gocui.Key) == gocui.KeyCtrlK {
			return "ctrl+k"
		}
		key = int(b.Key.(gocui.Key))
	}

	// special keys
	switch key {
	case 27:
		return "esc"
	case 13:
		return "enter"
	case 32:
		return "space"
	case 65514:
		return "►"
	case 65515:
		return "◄"
	case 65517:
		return "▲"
	case 65516:
		return "▼"
	case 65508:
		return "PgUp"
	case 65507:
		return "PgDn"
	case 9:
		return "tab"
	}

	return string(key)
}

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*Binding {
	bindings := []*Binding{
		{
			ViewName: "",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		}, {
			ViewName: "",
			Key:      'Q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuitWithoutChangingDirectory,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		}, {
			ViewName: "",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		}, {
			ViewName:    "",
			Key:         gocui.KeyPgup,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
		}, {
			ViewName:    "",
			Key:         gocui.KeyPgdn,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
		}, {
			ViewName: "",
			Key:      'K',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		}, {
			ViewName: "",
			Key:      'J',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlU,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlD,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		}, {
			ViewName:    "",
			Key:         'm',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.SLocalize("ViewMergeRebaseOptions"),
		}, {
			ViewName:    "",
			Key:         'P',
			Modifier:    gocui.ModNone,
			Handler:     gui.pushFiles,
			Description: gui.Tr.SLocalize("push"),
		}, {
			ViewName:    "",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePullFiles,
			Description: gui.Tr.SLocalize("pull"),
		}, {
			ViewName:    "",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefresh,
			Description: gui.Tr.SLocalize("refresh"),
		}, {
			ViewName: "",
			Key:      'x',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		}, {
			ViewName: "",
			Key:      '?',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		}, {
			ViewName: "",
			Key:      gocui.MouseMiddle,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlP,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreatePatchOptionsMenu,
		}, {
			ViewName:    "status",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.SLocalize("EditConfig"),
		}, {
			ViewName:    "status",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.SLocalize("OpenConfig"),
		}, {
			ViewName:    "status",
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.SLocalize("checkForUpdate"),
		}, {
			ViewName:    "status",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRecentReposMenu,
			Description: gui.Tr.SLocalize("SwitchRepo"),
		},
		{
			ViewName:    "files",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.SLocalize("CommitChanges"),
		},
		{
			ViewName:    "files",
			Key:         'w',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.SLocalize("commitChangesWithoutHook"),
		}, {
			ViewName:    "files",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAmendCommitPress,
			Description: gui.Tr.SLocalize("AmendLastCommit"),
		}, {
			ViewName:    "files",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.SLocalize("CommitChangesWithEditor"),
		}, {
			ViewName:    "files",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFilePress,
			Description: gui.Tr.SLocalize("toggleStaged"),
		}, {
			ViewName:    "files",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.Tr.SLocalize("viewDiscardOptions"),
		}, {
			ViewName:    "files",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		}, {
			ViewName:    "files",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		}, {
			ViewName:    "files",
			Key:         'i',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.SLocalize("ignoreFile"),
		}, {
			ViewName:    "files",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.SLocalize("refreshFiles"),
		}, {
			ViewName:    "files",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashChanges,
			Description: gui.Tr.SLocalize("stashAllChanges"),
		}, {
			ViewName:    "files",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateStashMenu,
			Description: gui.Tr.SLocalize("viewStashOptions"),
		}, {
			ViewName:    "files",
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageAll,
			Description: gui.Tr.SLocalize("toggleStagedAll"),
		}, {
			ViewName:    "files",
			Key:         'D',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateResetMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		}, {
			ViewName:    "files",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.SLocalize("StageLines"),
		}, {
			ViewName:    "files",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.SLocalize("fetch"),
		}, {
			ViewName:    "files",
			Key:         'X',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCustomCommand,
			Description: gui.Tr.SLocalize("executeCustomCommand"),
		}, {
			ViewName:    "branches",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleBranchPress,
			Description: gui.Tr.SLocalize("checkout"),
		}, {
			ViewName:    "branches",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.SLocalize("createPullRequest"),
		}, {
			ViewName:    "branches",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.SLocalize("checkoutByName"),
		}, {
			ViewName:    "branches",
			Key:         'F',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.SLocalize("forceCheckout"),
		}, {
			ViewName:    "branches",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleNewBranch,
			Description: gui.Tr.SLocalize("newBranch"),
		}, {
			ViewName:    "branches",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		}, {
			ViewName:    "branches",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRebase,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		}, {
			ViewName:    "branches",
			Key:         'M',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleMerge,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		}, {
			ViewName:    "branches",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFastForward,
			Description: gui.Tr.SLocalize("FastForward"),
		}, {
			ViewName:    "commits",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.SLocalize("squashDown"),
		}, {
			ViewName:    "commits",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommit,
			Description: gui.Tr.SLocalize("renameCommit"),
		}, {
			ViewName:    "commits",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommitEditor,
			Description: gui.Tr.SLocalize("renameCommitEditor"),
		}, {
			ViewName:    "commits",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateCommitResetMenu,
			Description: gui.Tr.SLocalize("resetToThisCommit"),
		}, {
			ViewName:    "commits",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.SLocalize("fixupCommit"),
		}, {
			ViewName:    "commits",
			Key:         'F',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateFixupCommit,
			Description: gui.Tr.SLocalize("createFixupCommit"),
		}, {
			ViewName:    "commits",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSquashAllAboveFixupCommits,
			Description: gui.Tr.SLocalize("squashAboveCommits"),
		}, {
			ViewName:    "commits",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.SLocalize("deleteCommit"),
		}, {
			ViewName:    "commits",
			Key:         gocui.KeyCtrlJ,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.SLocalize("moveDownCommit"),
		}, {
			ViewName:    "commits",
			Key:         gocui.KeyCtrlK,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.SLocalize("moveUpCommit"),
		}, {
			ViewName:    "commits",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.SLocalize("editCommit"),
		}, {
			ViewName:    "commits",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.SLocalize("amendToCommit"),
		}, {
			ViewName:    "commits",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.SLocalize("pickCommit"),
		}, {
			ViewName:    "commits",
			Key:         't',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.SLocalize("revertCommit"),
		}, {
			ViewName:    "commits",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.SLocalize("cherryPickCopy"),
		}, {
			ViewName:    "commits",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.SLocalize("cherryPickCopyRange"),
		}, {
			ViewName:    "commits",
			Key:         'v',
			Modifier:    gocui.ModNone,
			Handler:     gui.HandlePasteCommits,
			Description: gui.Tr.SLocalize("pasteCommits"),
		}, {
			ViewName:    "commits",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSwitchToCommitFilesPanel,
			Description: gui.Tr.SLocalize("viewCommitFiles"),
		}, {
			ViewName:    "commits",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleDiffCommit,
			Description: gui.Tr.SLocalize("CommitsDiff"),
		}, {
			ViewName:    "stash",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashApply,
			Description: gui.Tr.SLocalize("apply"),
		}, {
			ViewName:    "stash",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashPop,
			Description: gui.Tr.SLocalize("pop"),
		}, {
			ViewName:    "stash",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashDrop,
			Description: gui.Tr.SLocalize("drop"),
		}, {
			ViewName: "commitMessage",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitConfirm,
		}, {
			ViewName: "commitMessage",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitClose,
		}, {
			ViewName: "credentials",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSubmitCredential,
		}, {
			ViewName: "credentials",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCloseCredentialsView,
		}, {
			ViewName: "menu",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		}, {
			ViewName: "menu",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		}, {
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleDonate,
		}, {
			ViewName:    "commitFiles",
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSwitchToCommitsPanel,
			Description: gui.Tr.SLocalize("goBack"),
		}, {
			ViewName:    "commitFiles",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.Tr.SLocalize("checkoutCommitFile"),
		}, {
			ViewName:    "commitFiles",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.Tr.SLocalize("discardOldFileChange"),
		},
		{
			ViewName:    "commitFiles",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.Tr.SLocalize("toggleAddToPatch"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEnterCommitFile,
			Description: gui.Tr.SLocalize("enterFile"),
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpSecondary,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownSecondary,
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: gui.nextView},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	for i, viewName := range []string{"status", "files", "branches", "commits", "stash"} {
		bindings = append(bindings, &Binding{ViewName: "", Key: rune(i+1) + '0', Modifier: gocui.ModNone, Handler: gui.goToSideView(viewName)})
	}

	listPanelMap := map[string]struct {
		prevLine func(*gocui.Gui, *gocui.View) error
		nextLine func(*gocui.Gui, *gocui.View) error
		onClick  func(*gocui.Gui, *gocui.View) error
	}{
		"menu":        {prevLine: gui.handleMenuPrevLine, nextLine: gui.handleMenuNextLine, onClick: gui.handleMenuClick},
		"files":       {prevLine: gui.handleFilesPrevLine, nextLine: gui.handleFilesNextLine, onClick: gui.handleFilesClick},
		"branches":    {prevLine: gui.handleBranchesPrevLine, nextLine: gui.handleBranchesNextLine, onClick: gui.handleBranchesClick},
		"commits":     {prevLine: gui.handleCommitsPrevLine, nextLine: gui.handleCommitsNextLine, onClick: gui.handleCommitsClick},
		"stash":       {prevLine: gui.handleStashPrevLine, nextLine: gui.handleStashNextLine, onClick: gui.handleStashEntrySelect},
		"status":      {onClick: gui.handleStatusClick},
		"commitFiles": {prevLine: gui.handleCommitFilesPrevLine, nextLine: gui.handleCommitFilesNextLine, onClick: gui.handleCommitFilesClick},
	}

	for viewName, functions := range listPanelMap {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: functions.onClick},
		}...)
	}

	return bindings
}

// GetCurrentKeybindings gets the list of keybindings given the current context
func (gui *Gui) GetCurrentKeybindings() []*Binding {
	bindings := gui.GetInitialKeybindings()
	currentContext := gui.State.Context
	contextBindings := gui.GetContextMap()[currentContext]

	return append(bindings, contextBindings...)
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	if err := gui.setInitialContext(); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) GetContextMap() map[string][]*Binding {
	return map[string][]*Binding{
		"normal": {
			{
				ViewName: "secondary",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseDownSecondary,
			},
			{
				ViewName:    "main",
				Key:         gocui.MouseWheelDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.scrollDownMain,
				Description: gui.Tr.SLocalize("ScrollDown"),
				Alternative: "fn+up",
			}, {
				ViewName:    "main",
				Key:         gocui.MouseWheelUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.scrollUpMain,
				Description: gui.Tr.SLocalize("ScrollUp"),
				Alternative: "fn+down",
			}, {
				ViewName: "main",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseDownMain,
			},
		},
		"staging": {
			{
				ViewName: "secondary",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModNone,
				Handler:  gui.handleTogglePanelClick,
			},
			{
				ViewName:    "main",
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleStagingEscape,
				Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectPrevLine,
				Description: gui.Tr.SLocalize("PrevLine"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectNextLine,
				Description: gui.Tr.SLocalize("NextLine"),
			}, {
				ViewName: "main",
				Key:      'k',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevLine,
			}, {
				ViewName: "main",
				Key:      'j',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextLine,
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowLeft,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectPrevHunk,
				Description: gui.Tr.SLocalize("PrevHunk"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowRight,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectNextHunk,
				Description: gui.Tr.SLocalize("NextHunk"),
			}, {
				ViewName: "main",
				Key:      'h',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevHunk,
			}, {
				ViewName: "main",
				Key:      'l',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextHunk,
			}, {
				ViewName:    "main",
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleStageSelection,
				Description: gui.Tr.SLocalize("StageSelection"),
			}, {
				ViewName:    "main",
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleResetSelection,
				Description: gui.Tr.SLocalize("ResetSelection"),
			}, {
				ViewName:    "main",
				Key:         'v',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleToggleSelectRange,
				Description: gui.Tr.SLocalize("ToggleDragSelect"),
			}, {
				ViewName:    "main",
				Key:         'a',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleToggleSelectHunk,
				Description: gui.Tr.SLocalize("ToggleSelectHunk"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyTab,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleTogglePanel,
				Description: gui.Tr.SLocalize("TogglePanel"),
			}, {
				ViewName: "main",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseDown,
			}, {
				ViewName: "main",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModMotion,
				Handler:  gui.handleMouseDrag,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelUp,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseScrollUp,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelDown,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseScrollDown,
			},
		},
		"patch-building": {
			{
				ViewName:    "main",
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleEscapePatchBuildingPanel,
				Description: gui.Tr.SLocalize("ExitLineByLineMode"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectPrevLine,
				Description: gui.Tr.SLocalize("PrevLine"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectNextLine,
				Description: gui.Tr.SLocalize("NextLine"),
			}, {
				ViewName: "main",
				Key:      'k',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevLine,
			}, {
				ViewName: "main",
				Key:      'j',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextLine,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelUp,
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevLine,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelDown,
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextLine,
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowLeft,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectPrevHunk,
				Description: gui.Tr.SLocalize("PrevHunk"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowRight,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectNextHunk,
				Description: gui.Tr.SLocalize("NextHunk"),
			}, {
				ViewName: "main",
				Key:      'h',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevHunk,
			}, {
				ViewName: "main",
				Key:      'l',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextHunk,
			}, {
				ViewName:    "main",
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleAddSelectionToPatch,
				Description: gui.Tr.SLocalize("StageSelection"),
			}, {
				ViewName:    "main",
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleRemoveSelectionFromPatch,
				Description: gui.Tr.SLocalize("ResetSelection"),
			}, {
				ViewName:    "main",
				Key:         'v',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleToggleSelectRange,
				Description: gui.Tr.SLocalize("ToggleDragSelect"),
			}, {
				ViewName:    "main",
				Key:         'a',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleToggleSelectHunk,
				Description: gui.Tr.SLocalize("ToggleSelectHunk"),
			}, {
				ViewName: "main",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseDown,
			}, {
				ViewName: "main",
				Key:      gocui.MouseLeft,
				Modifier: gocui.ModMotion,
				Handler:  gui.handleMouseDrag,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelUp,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseScrollUp,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelDown,
				Modifier: gocui.ModNone,
				Handler:  gui.handleMouseScrollDown,
			},
		},
		"merging": {
			{
				ViewName:    "main",
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleEscapeMerge,
				Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.handlePickHunk,
				Description: gui.Tr.SLocalize("PickHunk"),
			}, {
				ViewName:    "main",
				Key:         'b',
				Modifier:    gocui.ModNone,
				Handler:     gui.handlePickBothHunks,
				Description: gui.Tr.SLocalize("PickBothHunks"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowLeft,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectPrevConflict,
				Description: gui.Tr.SLocalize("PrevConflict"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowRight,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectNextConflict,
				Description: gui.Tr.SLocalize("NextConflict"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectTop,
				Description: gui.Tr.SLocalize("SelectTop"),
			}, {
				ViewName:    "main",
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleSelectBottom,
				Description: gui.Tr.SLocalize("SelectBottom"),
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelUp,
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectTop,
			}, {
				ViewName: "main",
				Key:      gocui.MouseWheelDown,
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectBottom,
			}, {
				ViewName: "main",
				Key:      'h',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectPrevConflict,
			}, {
				ViewName: "main",
				Key:      'l',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectNextConflict,
			}, {
				ViewName: "main",
				Key:      'k',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectTop,
			}, {
				ViewName: "main",
				Key:      'j',
				Modifier: gocui.ModNone,
				Handler:  gui.handleSelectBottom,
			}, {
				ViewName:    "main",
				Key:         'z',
				Modifier:    gocui.ModNone,
				Handler:     gui.handlePopFileSnapshot,
				Description: gui.Tr.SLocalize("Undo"),
			}, {
				ViewName:    "main",
				Key:         'e',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleFileEdit,
				Description: gui.Tr.SLocalize("editFile"),
			}, {
				ViewName:    "main",
				Key:         'o',
				Modifier:    gocui.ModNone,
				Handler:     gui.handleFileOpen,
				Description: gui.Tr.SLocalize("openFile"),
			},
		},
	}
}
