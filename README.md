![Cascade](media/banner.png)
Ever had that sinking feeling when your system goes down at 3 AM? Yeah, me too. It sucks.

Here's the kicker: with distributed systems, failures aren't an exception, they're the norm. Our systems are complex beasts, and they're gonna break. Period. Here's what a typical microservice looks like.

![alt text](media/microservice.png)

But what if I told you that failure could be a feature and not a bug.

Enter Cascade. We break things on purpose so they don't break by accident. Simple as that!

Cascade doesn't just test your system. It pushes it to the brink, finds the weak spots, and helps you patch them up. All before your users even notice.

Think about it:
- What happens if your database suddenly decides to take a do a replication de-tour?
- If a critical service goes MIA, does your system curl up and cry, or does it power through?

These aren't hypotheticals. They're typical Tuesdays.

But we don't just break things and leave you to pick up the pieces. Oh no. We're not savages. (well, maybe a little).

Cascade lets you simulate these scenarios in a controlled environment. Crash and burn all you want - no real users harmed in the process.

Netflix uses chaos engineering to ensure your binge-watching is never interrupted.
Zomato uses it to keep your impulsive orders flowing 24/7.
Even NASA uses it. Yeah, freaking NASA boi.

If it's good enough for rocket scientists, it's prolly good enough for your app too.

Alright, going to get a bit nerdy for a sec. Chaos engineering isn't just about randomly unplugging servers for kicks (though that can be fun).

It's a discipline. A science. An art form.

Chaos Engineering Principles 101:

1. Start with a Hypothesis
   We're not savages. We don't just break things willy-nilly. We form a hypothesis about how our system should behave under stress. Then we test it. It's the scientific method, but with more explosions.

2. Minimize Blast Radius
   We're not trying to nuke your production environment. We start small, in controlled environments, and gradually increase the scope.

3. Run Experiments in Production
   "But that's crazy!" I hear you cry. Nope, that's confidence. Real chaos engineering means testing in the real world. Because your staging environment is a beautiful lie.

4. Automate Experiments
   Manual chaos is so 2010. We're talking continuous, automated chaos. Because if your chaos isn't continuous, neither is your resilience.

5. Measure Everything
   If you're not measuring, you're just breaking stuff for fun. We capture metrics on everything. Latency, traffic, errors, you name it. Data is king, and we're building an empire.

## Contributing

Feel free to contribute to Cascade by sending us your suggestions, bug
reports, or cat videos. Contributions are what make the open source community
such an amazing place to be learn, inspire, and create. Any contributions you
make are **greatly appreciated**.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.
Cascade is provided "as is" and comes with absolutely no guarantees.
If it breaks your system, well, that's kind of the point, isn't it? Congratulations, you're now doing chaos engineering!

Use at your own risk. Side effects may include improved system resilience, fewer 3 AM panic attacks, and an irresistible urge to push big red buttons.

## Credits

Created by engineer who wants to take down prod, loves to break things for a living and sleep soundly at night.

Special thanks to Murphy's Law for the constant inspiration.