# fsw-calendar

A small Go implementation for GCP Cloud Functions to automatically register the schedules of the running slots in Fuji International Speedway.

## Calendars

You can import the calendars using iCal URLs below. Only listed categories are supported now. There's no plan to support other categories.

- [SS-4](https://calendar.google.com/calendar/ical/25dfc5f523279ce53fbeb5a90eecfe45b2bdc8fdbe8a17bbae28ee77a7f53753%40group.calendar.google.com/public/basic.ics)
- [T-4](https://calendar.google.com/calendar/ical/0a31d519fe9246513b2bc7937d0191d60058eb3c74ec8dea8481642331c48305%40group.calendar.google.com/public/basic.ics)
- [NS-4](https://calendar.google.com/calendar/ical/4860bb7c2c62bf60cd5264495d3e739e06d82eaf66fdfb6f0e3f77e053d1b6fe%40group.calendar.google.com/public/basic.ics)

Each calendar is updated every early morning. If the original schedule has changed, the calendar follows it behind 24 hours at maximum.

## Disclaimer

The calendars might become unavailable without any notice in advance.
