package tcp_ip_proxys

import (
	"time"

	"github.com/alex60217101990/proxy.git/external/consts"
	"github.com/alex60217101990/proxy.git/external/logger"
	"github.com/alex60217101990/proxy.git/external/models"
	"github.com/jaredfolkins/badactor"
	"github.com/tevino/abool"
)

type FirewallBan struct {
	st          *badactor.Studio
	closed      *abool.AtomicBool
	ruleCh      chan models.BanObject
	banEventsCh chan models.BanEvent
}

func NewFirewallBan() Firewall {
	firewall := &FirewallBan{
		st:          badactor.NewStudio(consts.StudioCapacity),
		ruleCh:      make(chan models.BanObject, 10),
		banEventsCh: make(chan models.BanEvent, 200),
		closed:      abool.New(),
	}
	firewall.closed.UnSet()
	return firewall
}

func (f *FirewallBan) EmmitRule(rule *models.BanObject) {
	select {
	case f.ruleCh <- *rule:
	default:
		logger.Sugar.Warn("Emmit new firewall rule failed.")
	}
}

func (f *FirewallBan) EmmitBanEvent(event *models.BanEvent) {
	select {
	case f.banEventsCh <- *event:
	default:
		logger.Sugar.Warn("Emmit new firewall ban event failed.")
	}
}

func (f *FirewallBan) ruleLoop() {
	go func() {
		for !f.closed.IsSet() {
			select {
			case rule, ok := <-f.ruleCh:
				if ok {
					f.EmmitRule(&rule)
				}
			}
		}
	}()
}

func (f *FirewallBan) IsBan(ip string) bool {
	return f.st.IsJailed(ip)
}

func (f *FirewallBan) banEventLoop() {
	go func() {
		for !f.closed.IsSet() {
			select {
			case event, ok := <-f.banEventsCh:
				if ok {
					// action fails, increment infraction
					err := f.st.Infraction(event.IP, event.Name)
					if err != nil {
						logger.Sugar.Errorf("[%v] has err %v", event.IP, err)
					}
					// action fails, increment infraction
					i, err := f.st.Strikes(event.IP, event.Name)
					logger.Sugar.Errorf("[%v] has %v Strikes %v", event.IP, i, err)
				}
			}
		}
	}()
}

func (f *FirewallBan) LoadAllRules(rules []*models.BanObject) (err error) {
	for _, rule := range rules {
		// create and add rule
		newRule := &badactor.Rule{
			Name:        rule.Name,
			Message:     "You have failed to connect",
			StrikeLimit: int(rule.StrikeLimit),
			ExpireBase:  rule.ExpireBase,
			Sentence:    rule.Sentence,
		}
		// add the rule to the stack
		f.st.AddRule(newRule)
	}
	// creates the Directors who act as the Buckets in our sharding cache
	err = f.st.CreateDirectors(consts.StudioCapacity)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

func (f *FirewallBan) StudioLoop() {
	// Start the reaper
	f.st.StartReaper(time.Minute * time.Duration(60))
}
