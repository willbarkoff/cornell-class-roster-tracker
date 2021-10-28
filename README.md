# Cornell Class Roster Tracker

Sometimes when you're picking a class, you wonder: "Will this class be full by
the time I enroll?" Well, my plan is to collect data every hour and see when
classes fill.

There's probably a better way (each file is ~1.5MB * 24 collections / day * 7
days of preenroll) = approx 0.25GB of data, which I intended to display on a
website but am now thinking I could just generate a graph for every class
because nobody wants to load a 0.25GB website but idk man?

Right now, the plan is to run it on a cron job on `ugclinux`. We'll see how that
works out lol. It takes about 4 minutes to run, because the
[course roster](https://classes.cornell.edu) developers kindly request that you
don't make more than 1 request per second, so we wait a second between every
request.

This is going to be interesting and it might totally fail.

Also `ugclinux` doesn't have Go on it so you need to cross-compile (or if you're
on x86 linux i guess just compile) for it:

```
Linux en-ci-cisugcl16 5.4.0-80-generic #90-Ubuntu SMP Fri Jul 9 22:49:44 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
```

To cross compile:

```
GOOS=linux GOARCH=amd64
```

oh ALSO i might miss a few datapoints on ugclinux:

```
System maintenance: Tuesday or Thursday between 5am-7am these servers may be rebooted
```
