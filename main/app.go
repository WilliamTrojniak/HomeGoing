package main

import (
	"fmt"
	"gohome/dotmanager"
	"gohome/multiselect"
	"path"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
  Help key.Binding
  Quit key.Binding
  Refresh key.Binding
}

var keys = keyMap{
  Quit: key.NewBinding(
    key.WithKeys("q", "ctrl+c"),
    key.WithHelp("q", "quit")),
  Help: key.NewBinding(
    key.WithKeys("?"),
    key.WithHelp("?", "show help")),
  Refresh: key.NewBinding(
    key.WithKeys("r"),
    key.WithHelp("r", "refresh")),
}

func (k keyMap) ShortHelp() []key.Binding {
  return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
  return [][]key.Binding{{k.Help},{k.Quit}};
}


type app struct {
  help help.Model
  keys keyMap
  configFilePath string
  config *dotmanager.DotConfig
  linkSelector *multiselect.MultiSelectModel[dotmanager.SymLink]

  error
  isQuitting bool
}

func newApp(configFilePath string) *app {
  help := help.New();
  help.ShowAll = true;
  options := make([]dotmanager.SymLink, 0);

  return &app{
    help: help, 
    keys: keys, 
    configFilePath: path.Clean(configFilePath),
    linkSelector: multiselect.NewMultiSelect(options),
  };
}

func (m app) Init() tea.Cmd {
  return nil;
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    // TODO: Handle resizing
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.Quit):
      m.isQuitting = true;
      return m, tea.Quit
    case key.Matches(msg, m.keys.Help):
    // TODO Implement
    case key.Matches(msg, m.keys.Refresh):
      return m, getDotfilesConfig(m.configFilePath);
    case key.Matches(msg, key.NewBinding(key.WithKeys("l"))):
      if m.config != nil && len(m.config.GetModules()) > 0 {
        return m, linkModule(m.config.GetModules()[0]);
      }
    return m, nil;
    case key.Matches(msg, key.NewBinding(key.WithKeys("u"))):
      if m.config != nil && len(m.config.GetModules()) > 0 {
        return m, unlinkModule(m.config.GetModules()[0]);
      }
    return m, nil;
    }
  case GetDotfilesConfigMsg:
    m.config = msg.config;
  case GetActiveSymLinksMsg:
    m.linkSelector = multiselect.NewMultiSelect(msg.links);
    return m, nil;
  case LinkDotModuleMsg:
    return m, nil;
  case ErrMsg:
    m.error = msg.error;
   // TODO: Log other errors
    if msg.IsFatal() {
      return m, tea.Quit;
    }
  }
  

  return m, nil;
}

func (m app) View() string {
  if m.isQuitting {
    if m.error != nil {
      return fmt.Sprintf("A fatal error occurred: %v", m.error);
    }
    return "";
  }
  if m.error != nil {
    return fmt.Sprintf("%v", m.error);
  }
  var s string;
  s += "Dotfiles Source Directory: ";
  if m.config == nil {
    return s;
  }
  s += "Loaded";
  s += "\n";
  for _, mod := range m.config.GetModules() {
    status, _ := mod.GetLinkStatus();
    s += status.String() + " " + mod.GetName() + ": " +  mod.GetSrc() + " -> " + mod.GetDest() + "\n";
  }
  s += m.linkSelector.View();
  return s;
}
