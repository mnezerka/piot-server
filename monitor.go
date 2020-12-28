package main

import (
    "fmt"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "piot-server/config"
    "time"
)

type Monitor struct {
    log *logging.Logger
    db *mongo.Database
    mailClient IMailClient
    things *Things
    params *config.Parameters
    users *Users
}

func NewMonitor(log *logging.Logger,
                db *mongo.Database,
                mailClient IMailClient,
                things *Things,
                params *config.Parameters,
                users *Users) *Monitor {
                    return &Monitor{
                        log: log,
                        db: db,
                        mailClient: mailClient,
                        things: things,
                        params: params,
                        users: users}
}

func (m *Monitor) Check() {
    m.log.Infof("Monitor check started")

    var msgLines []string

    // get all enabled things
    filter := bson.M{"enabled": true}
    things, err := m.things.GetFiltered(filter)
    if err != nil {
        m.log.Errorf("Monitor check error, falied fetching of things: %s", err.Error())
        return
    }

    for i := 0; i < len(things); i++ {

        thing := things[i]

        // skip things where last seen interval is not set
        if thing.LastSeenInterval == 0 {
            continue
        }

        diff := int32(time.Now().Unix()) - thing.LastSeen

        if diff > thing.LastSeenInterval {

            lastSeen := time.Unix(int64(thing.LastSeen), 0)

            msgLine := fmt.Sprintf("%s (LastSeen: %s, LastSeenInterval: %d sec., Id: %s)",
                thing.Name,
                lastSeen,
                thing.LastSeenInterval,
                thing.Id.Hex())

            m.log.Warningf("Thing %s did not respond in defined interval", msgLine)
            msgLines = append(msgLines, msgLine)
        }
    }

    m.log.Infof("Monitor check - %d (out of %d) not responding things detected", len(msgLines), len(things))

    if len(msgLines) > 0 {

        msg := "Following things didn't respond in defined interval:\n\n"
        for i := 0; i < len(msgLines); i++ {
            msg += msgLines[i] + "\n"
        }

        // get admin users
        var adminEmails []string
        admins, err := m.users.GetAdmins()
        for i := 0; i < len(admins); i++ {
            adminEmails = append(adminEmails, admins[i].Email)
        }

        err = m.mailClient.SendMail("[piot][alarm] Not Available Devices", m.params.MailFrom, adminEmails, msg)
        if err != nil {
            m.log.Error(err)
        }
    }

    m.log.Infof("Monitor check finished")
}
