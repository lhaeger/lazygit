package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

var cyclableViews = []string{"status", "files", "branches", "commits", "stash"}

func (gui *Gui) refreshSidePanels(g *gocui.Gui) error {
	if err := gui.refreshBranches(g); err != nil {
		return err
	}
	if err := gui.refreshFiles(); err != nil {
		return err
	}
	if err := gui.refreshCommits(g); err != nil {
		return err
	}

	return gui.refreshStashEntries(g)
}

func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[len(cyclableViews)-1] {
		focusedViewName = cyclableViews[0]
	} else {
		// if we're in the commitFiles view we'll act like we're in the commits view
		viewName := v.Name()
		if viewName == "commitFiles" {
			viewName = "commits"
		}
		for i := range cyclableViews {
			if viewName == cyclableViews[i] {
				focusedViewName = cyclableViews[i+1]
				break
			}
			if i == len(cyclableViews)-1 {
				message := gui.Tr.TemplateLocalize(
					"IssntListOfViews",
					Teml{
						"name": viewName,
					},
				)
				gui.Log.Info(message)
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
	return gui.switchFocus(g, v, focusedView)
}

func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[0] {
		focusedViewName = cyclableViews[len(cyclableViews)-1]
	} else {
		// if we're in the commitFiles view we'll act like we're in the commits view
		viewName := v.Name()
		if viewName == "commitFiles" {
			viewName = "commits"
		}
		for i := range cyclableViews {
			if viewName == cyclableViews[i] {
				focusedViewName = cyclableViews[i-1] // TODO: make this work properly
				break
			}
			if i == len(cyclableViews)-1 {
				message := gui.Tr.TemplateLocalize(
					"IssntListOfViews",
					Teml{
						"name": viewName,
					},
				)
				gui.Log.Info(message)
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
	return gui.switchFocus(g, v, focusedView)
}

func (gui *Gui) newLineFocused(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "menu":
		return gui.handleMenuSelect(g, v)
	case "status":
		return gui.handleStatusSelect(g, v)
	case "files":
		return gui.handleFileSelect(g, v, false)
	case "branches":
		return gui.handleBranchSelect(g, v)
	case "commits":
		return gui.handleCommitSelect(g, v)
	case "commitFiles":
		return gui.handleCommitFileSelect(g, v)
	case "stash":
		return gui.handleStashEntrySelect(g, v)
	case "confirmation":
		return nil
	case "commitMessage":
		return gui.handleCommitFocused(g, v)
	case "credentials":
		return gui.handleCredentialsViewFocused(g, v)
	case "main":
		if gui.State.Context == "merging" {
			return gui.refreshMergePanel()
		}
		v.Highlight = false
		return nil
	default:
		panic(gui.Tr.SLocalize("NoViewMachingNewLineFocusedSwitchStatement"))
	}
}

func (gui *Gui) returnFocus(g *gocui.Gui, v *gocui.View) error {
	previousView, err := g.View(gui.State.PreviousView)
	if err != nil {
		// always fall back to files view if there's no 'previous' view stored
		previousView, err = g.View("files")
		if err != nil {
			gui.Log.Error(err)
		}
	}
	return gui.switchFocus(g, v, previousView)
}

func (gui *Gui) goToSideView(sideViewName string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		view, err := g.View(sideViewName)
		if err != nil {
			gui.Log.Error(err)
			return nil
		}
		err = gui.closePopupPanels()
		if err != nil {
			gui.Log.Error(err)
			return nil
		}
		return gui.switchFocus(g, nil, view)
	}
}

func (gui *Gui) closePopupPanels() error {
	gui.onNewPopupPanel()
	err := gui.closeConfirmationPrompt(gui.g, true)
	if err != nil {
		gui.Log.Error(err)
		return err
	}
	return nil
}

// pass in oldView = nil if you don't want to be able to return to your old view
// TODO: move some of this logic into our onFocusLost and onFocus hooks
func (gui *Gui) switchFocus(g *gocui.Gui, oldView, newView *gocui.View) error {
	// we assume we'll never want to return focus to a popup panel i.e.
	// we should never stack popup panels
	if oldView != nil && !gui.isPopupPanel(oldView.Name()) {
		gui.State.PreviousView = oldView.Name()
	}

	gui.Log.Info("setting highlight to true for view" + newView.Name())
	message := gui.Tr.TemplateLocalize(
		"newFocusedViewIs",
		Teml{
			"newFocusedView": newView.Name(),
		},
	)
	gui.Log.Info(message)
	if _, err := g.SetCurrentView(newView.Name()); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop(newView.Name()); err != nil {
		return err
	}

	g.Cursor = newView.Editable

	if err := gui.renderPanelOptions(); err != nil {
		return err
	}

	return gui.newLineFocused(g, newView)
}

func (gui *Gui) resetOrigin(v *gocui.View) error {
	_ = v.SetCursor(0, 0)
	return v.SetOrigin(0, 0)
}

// if the cursor down past the last item, move it to the last line
func (gui *Gui) focusPoint(cx int, cy int, lineCount int, v *gocui.View) error {
	if cy < 0 || cy > lineCount {
		return nil
	}
	ox, oy := v.Origin()
	_, height := v.Size()

	ly := height - 1
	if ly == -1 {
		ly = 0
	}

	// if line is above origin, move origin and set cursor to zero
	// if line is below origin + height, move origin and set cursor to max
	// otherwise set cursor to value - origin
	if ly > lineCount {
		_ = v.SetCursor(cx, cy)
		_ = v.SetOrigin(ox, 0)
	} else if cy < oy {
		_ = v.SetCursor(cx, 0)
		_ = v.SetOrigin(ox, cy)
	} else if cy > oy+ly {
		_ = v.SetCursor(cx, ly)
		_ = v.SetOrigin(ox, cy-ly)
	} else {
		_ = v.SetCursor(cx, cy-oy)
	}
	return nil
}

func (gui *Gui) cleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return utils.NormalizeLinefeeds(output)
}

func (gui *Gui) setViewContent(g *gocui.Gui, v *gocui.View, s string) error {
	v.Clear()
	fmt.Fprint(v, gui.cleanString(s))
	return nil
}

// renderString resets the origin of a view and sets its content
func (gui *Gui) renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return nil // return gracefully if view has been deleted
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		return gui.setViewContent(gui.g, v, s)
	})
	return nil
}

func (gui *Gui) optionsMapToString(optionsMap map[string]string) string {
	optionsArray := make([]string, 0)
	for key, description := range optionsMap {
		optionsArray = append(optionsArray, key+": "+description)
	}
	sort.Strings(optionsArray)
	return strings.Join(optionsArray, ", ")
}

func (gui *Gui) renderOptionsMap(optionsMap map[string]string) error {
	return gui.renderString(gui.g, "options", gui.optionsMapToString(optionsMap))
}

// TODO: refactor properly
// i'm so sorry but had to add this getBranchesView
func (gui *Gui) getFilesView() *gocui.View {
	v, _ := gui.g.View("files")
	return v
}

func (gui *Gui) getCommitsView() *gocui.View {
	v, _ := gui.g.View("commits")
	return v
}

func (gui *Gui) getCommitMessageView() *gocui.View {
	v, _ := gui.g.View("commitMessage")
	return v
}

func (gui *Gui) getBranchesView() *gocui.View {
	v, _ := gui.g.View("branches")
	return v
}

func (gui *Gui) getMainView() *gocui.View {
	v, _ := gui.g.View("main")
	return v
}

func (gui *Gui) getSecondaryView() *gocui.View {
	v, _ := gui.g.View("secondary")
	return v
}

func (gui *Gui) getStashView() *gocui.View {
	v, _ := gui.g.View("stash")
	return v
}

func (gui *Gui) getCommitFilesView() *gocui.View {
	v, _ := gui.g.View("commitFiles")
	return v
}

func (gui *Gui) trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
}

func (gui *Gui) currentViewName() string {
	currentView := gui.g.CurrentView()
	return currentView.Name()
}

func (gui *Gui) resizeCurrentPopupPanel(g *gocui.Gui) error {
	v := g.CurrentView()
	if gui.isPopupPanel(v.Name()) {
		return gui.resizePopupPanel(g, v)
	}
	return nil
}

func (gui *Gui) resizePopupPanel(g *gocui.Gui, v *gocui.View) error {
	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	content := v.Buffer()
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, v.Wrap, content)
	vx0, vy0, vx1, vy1 := v.Dimensions()
	if vx0 == x0 && vy0 == y0 && vx1 == x1 && vy1 == y1 {
		return nil
	}
	gui.Log.Info(gui.Tr.SLocalize("resizingPopupPanel"))
	_, err := g.SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

// generalFocusLine takes a lineNumber to focus, and a bottomLine to ensure we can see
func (gui *Gui) generalFocusLine(lineNumber int, bottomLine int, v *gocui.View) error {
	_, height := v.Size()
	overScroll := bottomLine - height + 1
	if overScroll < 0 {
		overScroll = 0
	}
	if err := v.SetOrigin(0, overScroll); err != nil {
		return err
	}
	if err := v.SetCursor(0, lineNumber-overScroll); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) changeSelectedLine(line *int, total int, up bool) {
	if up {
		if *line == -1 || *line == 0 {
			return
		}

		*line--
	} else {
		if *line == -1 || *line == total-1 {
			return
		}

		*line++
	}
}

func (gui *Gui) refreshSelectedLine(line *int, total int) {
	if *line == -1 && total > 0 {
		*line = 0
	} else if total-1 < *line {
		*line = total - 1
	}
}

func (gui *Gui) renderListPanel(v *gocui.View, items interface{}) error {
	gui.g.Update(func(g *gocui.Gui) error {
		isFocused := gui.g.CurrentView().Name() == v.Name()
		list, err := utils.RenderList(items, isFocused)
		if err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		v.Clear()
		fmt.Fprint(v, list)
		return nil
	})
	return nil
}

func (gui *Gui) renderPanelOptions() error {
	currentView := gui.g.CurrentView()
	switch currentView.Name() {
	case "menu":
		return gui.renderMenuOptions()
	case "main":
		if gui.State.Context == "merging" {
			return gui.renderMergeOptions()
		}
	}
	return gui.renderGlobalOptions()
}

func (gui *Gui) handleFocusView(g *gocui.Gui, v *gocui.View) error {
	_, err := gui.g.SetCurrentView(v.Name())
	return err
}

func (gui *Gui) isPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "credentials" || viewName == "confirmation" || viewName == "menu"
}

func (gui *Gui) popupPanelFocused() bool {
	return gui.isPopupPanel(gui.currentViewName())
}

func (gui *Gui) handleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func(*gocui.Gui, *gocui.View) error) error {
	if gui.popupPanelFocused() && v != nil && !gui.isPopupPanel(v.Name()) {
		return nil
	}

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	newSelectedLine := v.SelectedLineIdx()

	if newSelectedLine < 0 {
		newSelectedLine = 0
	}

	if newSelectedLine > itemCount-1 {
		newSelectedLine = itemCount - 1
	}

	*selectedLine = newSelectedLine

	return handleSelect(gui.g, v)
}
