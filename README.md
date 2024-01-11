# fsw-calendar
[![test](https://github.com/yagihash/fsw-calendar/actions/workflows/test.yml/badge.svg)](https://github.com/yagihash/fsw-calendar/actions/workflows/test.yml)
[![Static Badge](https://img.shields.io/badge/test_coverage-blue)](https://yagihash.github.io/fsw-calendar/)

A small Go implementation for GCP Cloud Functions to automatically register the schedules of the membership driving slots at Fuji International Speedway.

## Calendars
> [!WARNING]
> The calendars are based on the official calendar which is provided by FSW.
> This is an unofficial use case and there's no guarantee of its availability.
> Thus the calendars might become unavailable without any notice in advance.

You can import the calendars using iCal URLs below. Only listed categories are supported now. There's no plan to support other categories.

- [SS-4](https://calendar.google.com/calendar/ical/25dfc5f523279ce53fbeb5a90eecfe45b2bdc8fdbe8a17bbae28ee77a7f53753%40group.calendar.google.com/public/basic.ics)
- [T-4](https://calendar.google.com/calendar/ical/0a31d519fe9246513b2bc7937d0191d60058eb3c74ec8dea8481642331c48305%40group.calendar.google.com/public/basic.ics)
- [NS-4](https://calendar.google.com/calendar/ical/4860bb7c2c62bf60cd5264495d3e739e06d82eaf66fdfb6f0e3f77e053d1b6fe%40group.calendar.google.com/public/basic.ics)
- [S-4](https://calendar.google.com/calendar/ical/ccc748b5eb57fc6d8d01d2a380f7833225b347a54b816d23a5dbabb209f263b3%40group.calendar.google.com/public/basic.ics)

Each calendar is updated every early morning. If the original schedule has changed, the calendar follows it behind 24 hours at maximum.

The calendars are also available below.

- [SS-4](https://calendar.google.com/calendar/embed?src=25dfc5f523279ce53fbeb5a90eecfe45b2bdc8fdbe8a17bbae28ee77a7f53753%40group.calendar.google.com&ctz=Asia%2FTokyo)
- [T-4](https://calendar.google.com/calendar/embed?src=0a31d519fe9246513b2bc7937d0191d60058eb3c74ec8dea8481642331c48305%40group.calendar.google.com&ctz=Asia%2FTokyo)
- [NS-4](https://calendar.google.com/calendar/embed?src=4860bb7c2c62bf60cd5264495d3e739e06d82eaf66fdfb6f0e3f77e053d1b6fe%40group.calendar.google.com&ctz=Asia%2FTokyo)
- [S-4](https://calendar.google.com/calendar/embed?src=ccc748b5eb57fc6d8d01d2a380f7833225b347a54b816d23a5dbabb209f263b3%40group.calendar.google.com&ctz=Asia%2FTokyo)
