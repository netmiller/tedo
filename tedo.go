package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	// "time"
	// "regexp"

	"./taskmanager"
	"github.com/xlab/tablewriter"		// alternate tablewriter
	"github.com/segmentio/go-prompt"

)

const usage = `
    Name:          tedo (Terminal Todo)
    -------        ----------------------
    $ tedo         Show all tasks
    $ tedo p       Show all pending tasks
    $ tedo s ID    Show detail view task of ID
    $ tedo a       Add a new task
    $ tedo m ID    Modify a task
    $ tedo rm ID   Remove task of ID from list
    $ tedo del     Remove latest task from list
    $ tedo c ID    Mark task of ID as completed
    $ tedo p ID    Mark task of ID as pending
    $ tedo flush   Flush the database!
`
//     $ tedo remind  Will send you a desktop notification
//     $ tedo service-start  Run task as service if you are using reminder
//     $ tedo service-stop   Unregister Task from service!
// `

const (
	// completedSign  = "\u2713"
	completedSign  = "ok"
	// pendingSign = "\u2613"
	// pendingSign = "\u2716"
	// pendingSign = "\u24c5"
	// pendingSign = "\u2757"
	// pendingSign = "\U0001F53B"
	pendingSign    = "\u24c5"
	dateTimeLayout = "2006-01-02 15:04"
	refreshRate    = 40
)

var (
	//task manager instance
	tm = taskmanager.New()
)

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	// fmt.Println("huu1")
	flag.Parse()
	// fmt.Println("huu2")

	cmd, args, argsLen := flag.Arg(0), flag.Args(), len(flag.Args())

	// fmt.Println("args:",args)

	switch {

	// näytetään kaikki
	case cmd == "" || cmd == "l" || cmd == "ls" && argsLen == 1:
		showTasksInTable(tm.GetAllTasks())

	// lisätään uusi task
	case cmd == "a" || cmd == "add" && argsLen >= 1:
		if len(args[1:]) <= 0 {
			fmt.Println(" Task description can not be empty \n")
			return
		}
		tm.Add(strings.Join(args[1:], " "), "", "")
		// successText(" Added to list: " + strings.Join(args[1:], " ") + " ")
		fmt.Println(" Added : " + strings.Join(args[1:], " ") + " ")

	case cmd == "p" || cmd == "pending" && argsLen == 1:
		showTasksInTable(tm.GetPendingTasks())

	case cmd == "del" || cmd == "delete" && argsLen == 1:
		p := prompt.Choose("Delete latest task ?", []string{"yes", "no"})
		if p == 1 {
			// warningText(" Task delete aboarted! ")
			fmt.Println(" Task delete canceled! ")
			return
		}
		err := tm.RemoveTask(tm.GetLastId())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// successText(" Removed latest task ")
		fmt.Println(" Removed latest task ")

	// poistetaan task
	case cmd == "r" || cmd == "rm" && argsLen == 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		// p := prompt.Choose("Do you want to delete task of id " + flag.Arg(1) + " ?", []string{"yes", "no"})
		p := prompt.Confirm("Delete task id " + flag.Arg(1) + " (y/n) >")
		// if p == 1 {
		if p != true {
			// warningText(" Task delete aboarted! ")
			fmt.Println(" Task delete canceled! ")
			return
		}
		err := tm.RemoveTask(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// successText(" Task " + strconv.Itoa(id) + " removed! ")
		fmt.Println(" Task " + strconv.Itoa(id) + " removed! ")

	case cmd == "e" || cmd == "m" || cmd == "u" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		ok, _ := tm.UpdateTask(id, strings.Join(args[2:], " "))
		// successText(ok)
		fmt.Println(ok)

	case cmd == "c" || cmd == "d" || cmd == "done" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.MarkAsCompleteTask(id)
		if err != nil {
			// errorText(err.Error())
			fmt.Println(err.Error())
			return
		}
		// successText(" " + completedSign + " " + task.Description)
		fmt.Println(" " + completedSign + " " + task.Description)

	case cmd == "i" || cmd == "p" || cmd == "pending" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.MarkAsPendingTask(id)
		if err != nil {
			// errorText(err.Error())
			fmt.Println(err.Error())
			return
		}
		// successText(" " + pendingMark() + " " + task.Description)
		fmt.Println(" " + pendingMark() + " " + task.Description)

	case cmd == "s" && argsLen == 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.GetTask(id)
		if err != nil {
			// errorText(err.Error())
			fmt.Println(err.Error())
			return
		}
		showTask(task)

	case cmd == "flush":
		p := prompt.Choose("Do you want to delete all tasks?", []string{"yes", "no"})
		if p == 1 {
			fmt.Println(" Flush aborted! ")
			return
		}
		err := tm.FlushDB()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// successText(" Database flushed successfully! ")
		fmt.Println(" Database flushed successfully! ")

	// case cmd == "reminder" || cmd == "remind" || cmd == "remind-me" && argsLen >= 1:
	// 	if len(args[1:]) <= 0 {
	// 		warningText(" Task/Reminder description can not be empty \n")
	// 		return
	// 	}
	// 	reminder := strings.Join(args[1:], " ")
	// 	action, actionWhen := parseReminder(reminder)
	// 	tm.Add(action, "", actionWhen)
	// 	// successText(" Reminder Added: " + action + " ")
	// 	fmt.Println(" Reminder Added: " + action + " ")

	// case cmd == "service-start" && argsLen == 1:
	// 	serviceStart()

	// case cmd == "service-force-start" && argsLen == 1:
	// 	serviceForceStart()

	// case cmd == "service-stop" && argsLen == 1:
	// 	serviceStop()

	// case cmd == "listen-reminder-queue" && argsLen == 1:
	// 	listenReminderQueue()

	case cmd == "h" || cmd == "v":
		// fmt.Fprint(os.Stderr, usage)
		fmt.Println(usage)
		return

	default:
		fmt.Println(" [No command found by " + cmd + "] ")
		// fmt.Fprint(os.Stderr, "\n"+usage)
		fmt.Println(usage)
		return
	}

}

// ----
func showTasksInTable(tasks taskmanager.Tasks) {
	fmt.Fprintln(os.Stdout, "")
	table := tablewriter.CreateTable()
	tablewriter.EnableUTF8()
	// table.AddHeaders("Id", "Description", completedSign+"/"+pendingMark(), "Created")
	table.AddHeaders("Id", "Description", "\u214f", "Created")

	for _, task := range tasks {
		//set completed icon
		status := pendingSign
		if task.Completed != "" {
			status = completedSign
		} else {
			status = pendingMark()
		}
		desc := task.Description
		if len(task.Description) > 57 {
			desc = task.Description[0:57] + " ..."
		}
		table.AddRow(strconv.Itoa(task.Id),	desc, status, task.Created )
	}
	fmt.Println(table.Render())
	fmt.Println("")
}

// ---
func showTask(task taskmanager.Task) {
	fmt.Fprintln(os.Stdout, "")
	table := tablewriter.CreateTable()
	tablewriter.EnableUTF8()
	table.AddHeaders("Key", "Value")
	// table.AddRow("Id", "["+strconv.Itoa(task.Id)+"]")
	w2 := strings.Fields(task.Description)
	t2 := "["+strconv.Itoa(task.Id)+"]"
	r2 := ""
	for _,r := range w2 {
		// fmt.Println(r)
		if (len(r2)+len(r)) > 85 {
			table.AddRow(t2, r2)
			t2 = " ..."
			r2 = r + " "
		} else {
			r2 += r + " "
		}
	}
	if len(r2) > 0 {
		table.AddRow(t2, r2)
	}
	table.AddRow(" Tag", task.Tag)
	table.AddRow(" Cre", task.Created)
	table.AddRow(" Upd", task.Updated)
	table.AddRow(" Uid", task.UID)
	fmt.Println(table.Render())
	fmt.Println("")
}


// ---
func printText(str string) {
	fmt.Fprintf(os.Stdout, str+"\n")
}

// ---
func pendingMark() string {
	pending := pendingSign
	if runtime.GOOS == "windows" {
		pending = "x"
	}
	return pending
}

//parse reminder
// func parseReminder(reminder string) (string, string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			errorText(" Your reminder does not contain any date time reference! ")
// 			os.Exit(1)
// 		}
// 	}()
// 	w := when.New(nil)
// 	w.Add(en.All...)
// 	w.Add(common.All...)
// 	r, _ := w.Parse(reminder, time.Now())
// 	action := strings.Replace(reminder, reminder[r.Index:r.Index+len(r.Text)], "", -1)
// 	actionTime := r.Time.Format(dateTimeLayout)
// 	return action, actionTime
// }

//listen for reminder queue
// func listenReminderQueue() {
// 	for {
// 		rm := taskmanager.New()
// 		reminderList := rm.GetReminderTasks()
// 		now := time.Now().Format(dateTimeLayout)
// 		for _, r := range reminderList {
// 			if r.RemindAt == now {
// 				desktopNotifier("Task Reminder!", r.Description)
// 				rm.MarkAsCompleteTask(r.Id)
// 			}
// 		}
// 		time.Sleep(time.Second * refreshRate)
// 	}
// }

//send desktop notification
// func desktopNotifier(title, body string) {
// 	notify = notificator.New(notificator.Options{
// 		DefaultIcon: "default.png",
// 		AppName:     "Terminal Task",
// 	})
// 	notify.Push(title, body, "default.png", notificator.UR_NORMAL)
// }

//enable auto start
// func serviceStart() {
// 	if service.IsEnabled() {
// 		warningText("Task is already enabled as service!")
// 	} else {
// 		if err := service.Enable(); err != nil {
// 			errorText(err.Error())
// 		}
// 		successText("Task has been registered as service!")
// 	}
// }

//disable auto start
// func serviceStop() {
// 	if service.IsEnabled() {
// 		if err := service.Disable(); err != nil {
// 			errorText(err.Error())
// 		}
// 		successText("Task has been removed from service!")
// 	} else {
// 		warningText("Task was not registered as service!")
// 	}
// }

//force start
// func serviceForceStart() {
// 	serviceStop()
// 	serviceStart()
// }
