package loading

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var Data Config

func Initialize() {
	Data = *Load()
}

func Load() *Config {
	return &Config{
		Proxies:         *loadProxies(),
		Accounts:        *loadAccounts(),
		Profiles:        *loadProfiles(),
		Tasks:           *loadTasks(),
		QuicktaskGroups: *loadQuicktasks(),
		Settings:        *loadSettings(),
	}
}

func loadSettings() *Settings {
	jsonFile, err := os.Open("../settings.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var settings Settings

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &settings.Settings)
	if err != nil {
		return nil
	}
	return &settings
}

func loadProxies() *Proxies {
	mainGroupUUID := proxy.CreateProxyGroup("main")

	f, err := os.Open("../proxies.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		proxyUUID := proxy.CreateProxy(scanner.Text())
		err := proxy.SetProxyToProxyGroup(proxyUUID, mainGroupUUID)
		if err != nil {
			log.Fatalf("error setting proxy %s to proxy group %s", proxyUUID, mainGroupUUID)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	mainGroupObject, _ := proxy.GetProxyGroup(mainGroupUUID)
	return &Proxies{
		Proxies: []proxy.ProxyGroup{
			*mainGroupObject,
		},
	}
}

func loadProfiles() *Profiles {
	f, err := os.Open("../profiles.csv")
	if err != nil {
		log.Fatal(err)
	}

	var profiles Profiles

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	c := 0
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if c == 0 {
			c += 1
			continue
		}

		profileId := profile.CreateProfile(&profile.Profile{
			Name: rec[0],
			Address: struct {
				Name        string `json:"name"`
				Email       string `json:"email"`
				HomePhone   string `json:"homePhone"`
				MobilePhone string `json:"mobilePhone"`
				Address     string `json:"address"`
				ZipCode     string `json:"zipCode"`
				City        string `json:"city"`
				Area        string `json:"area"`
				Prefecture  string `json:"prefecture"`
			}{
				Name:        rec[1] + " " + rec[2],
				Email:       rec[3],
				HomePhone:   rec[9],
				MobilePhone: rec[10],
				Address:     rec[4],
				ZipCode:     rec[5],
				City:        rec[6],
				Area:        rec[7],
				Prefecture:  rec[8],
			},
		})

		profileObject, _ := profile.GetProfileById(profileId)
		profiles.Profiles = append(profiles.Profiles, *profileObject)
	}

	return &profiles
}

func loadAccounts() *Accounts {
	paths := []string{
		"../sites/athletesfoot/accounts_athletesfoot.csv",
		"../sites/fuel/accounts_fuel.csv",
		"../sites/slamdunk/accounts_slamdunk.csv",
		"../sites/buzzsneakers/accounts_buzzsneakers.csv",
		"../sites/europesports/accounts_europesports.csv",
	}

	var accounts Accounts
	accounts.Accounts = make(map[int][]account.Account)

	for siteId, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		c := 0
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if c == 0 {
				c += 1
				continue
			}

			if rec[0] == "" || rec[1] == "" {
				continue
			}

			account.CreateAccount(&account.Account{
				SiteId:   siteId,
				Email:    rec[0],
				Password: rec[1],
			})
			accountObject, _ := account.GetAccount(siteId, rec[0])
			accounts.Accounts[siteId] = append(accounts.Accounts[siteId], *accountObject)

		}
		f.Close()
	}

	return &accounts
}

func loadQuicktasks() *map[int][]QuicktaskGroup {
	paths := []string{
		"../sites/athletesfoot/quicktasks_athletesfoot.csv",
		"../sites/fuel/quicktasks_fuel.csv",
		"../sites/slamdunk/quicktasks_slamdunk.csv",
		"../sites/buzzsneakers/quicktasks_buzzsneakers.csv",
		"../sites/europesports/quicktasks_europesports.csv",
	}

	quicktasks := make(map[int][]QuicktaskGroup)

	for siteId, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		// read csv values using csv.Reader
		QTGroup := make([]QuicktaskGroup, 0)

		csvReader := csv.NewReader(f)
		c := 0
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if c == 0 {
				c += 1
				continue
			}

			if rec[0] == "" {
				continue
			}

			QTGroup = append(QTGroup, QuicktaskGroup{
				ProfileName:  rec[0],
				AccountEmail: rec[1],
			})
		}

		quicktasks[siteId] = QTGroup

		f.Close()
	}

	return &quicktasks
}

func loadTasks() *Tasks {
	paths := []string{
		"../sites/athletesfoot/tasks_athletesfoot.csv",
		"../sites/fuel/tasks_fuel.csv",
		"../sites/slamdunk/tasks_slamdunk.csv",
		"../sites/buzzsneakers/tasks_buzzsneakers.csv",
		"../sites/europesports/tasks_europesports.csv",
	}

	var tasks Tasks
	tasks.Tasks = make(map[int][]string)

	for siteId, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		c := 0
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if c == 0 {
				c += 1
				continue
			}

			if rec[0] == "" {
				continue
			}

			profileId, _ := profile.GetProfileIDByName(rec[0])
			proxyGroupId, _ := proxy.GetProxyGroupIDByName("main")
			delay, err := strconv.Atoi(rec[7])
			if err != nil {
				log.Fatal("Failed to convert delay in a task to int.")
			}
			taskQuantity, err := strconv.Atoi(rec[6])
			if err != nil {
				log.Fatal("Failed to convert quantity in a task to int.")
			}

			for i := 0; i < taskQuantity; i++ {
				taskUUID := task.CreateTask(
					TaskModeAndSiteIDToRegisteredSiteName[rec[5]+","+strconv.Itoa(siteId)], // registered site name
					rec[2], // product info
					strings.TrimSpace(strings.ToLower(rec[3])), // size
					profileId,    // profile id
					proxyGroupId, // proxy group id
					rec[1],       // account info
					strings.TrimSpace(strings.ToLower(rec[5])), // task type
					strings.TrimSpace(strings.ToLower(rec[4])), // task mode
					delay, // delay
				)

				tasks.Tasks[siteId] = append(tasks.Tasks[siteId], taskUUID)
			}
		}

		f.Close()
	}

	return &tasks
}
