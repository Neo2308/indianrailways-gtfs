package server

func getServiceId(runningDays string) string {
	serviceId := ""
	if runningDays[0] == '1' {
		serviceId += "Mon-"
	}
	if runningDays[1] == '1' {
		serviceId += "Tue-"
	}
	if runningDays[2] == '1' {
		serviceId += "Wed-"
	}
	if runningDays[3] == '1' {
		serviceId += "Thu-"
	}
	if runningDays[4] == '1' {
		serviceId += "Fri-"
	}
	if runningDays[5] == '1' {
		serviceId += "Sat-"
	}
	if runningDays[6] == '1' {
		serviceId += "Sun-"
	}
	return serviceId[:len(serviceId)-1]
}
