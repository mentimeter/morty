<p align="center">
  <img alt="logo" src="https://user-images.githubusercontent.com/7828615/95340563-9b261880-08b5-11eb-86f7-156d6c29b44d.png" height="150" />
  <h3 align="center">Morty</h3>
  <p align="center">An action that turns a GitHub Repository into an organized collection of post-mortems</p>
</p>

---

> "There is no better way to learn than to document what has broken in the past. History is about learning from everyone’s mistakes.
> Be thorough, be honest, but most of all, ask hard questions. Look for specific actions that might prevent such an outage from recurring,
> not just tactically, but also strategically. Ensure that everyone within the company can learn what you have learned by publishing and organizing postmortems."

~ _From the "Emergency Response" Chapter in "Site Reliability Engineering: How Google Runs Production Systems"_

**Check out the [example repository](https://github.com/mentimeter/example-post-mortems) to see what it's like!**

Morty is a GitHub Action that gives you an overview of your post-mortems. It parses your post-mortems _written in markdown_, and gives you some friendly advice on running good post-mortems on the way 📈

![morty-example](https://user-images.githubusercontent.com/7828615/95354012-3b833980-08c4-11eb-83c0-845a8523d036.png)

## Getting started

Install the action in a (new) repository. You can create one from a [template repository](https://github.com/mentimeter/example-post-mortems) if you'd like!

To install the action, add a new workflow file `.github/workflows/morty.yml`. It should contain something like this:

```
name: Morty
on: [push, pull_request]

jobs:
  morty:
    runs-on: ubuntu-latest

    steps:
    - name: Organize mortems
      uses: mentimeter/morty@v1
      with:
        token: ${{ github.token }}
```

The easiest way to make a new post-mortem is to make a copy of the template that morty makes for you `post-mortems/template.md`.
There are also some more instructions there to help you get going.

## Why morty / post mortems as a repository?

- An incident history you can analyse
- Zero-overhead organization
- Easy to do global search
- Easy to follow normal pull request flow & collaborate
