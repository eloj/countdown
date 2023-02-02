
# Countdown Word Finder

[![Build status](https://github.com/eloj/countdown/actions/workflows/go.yml/badge.svg)](https://github.com/eloj/countdown/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/eloj/countdown)](https://goreportcard.com/report/github.com/eloj/countdown)

This is a Go demo project, which solves the problem: _"given a string of letters, which
words can be formed by using each letter from the string at most once."_

This word-game is played on the British TV panel show "Countdown" and its sister show
"8 Out Of 10 Cats... Does Countdown", hence the title of the repo. On the show, the
contestants get to pick if they want a random vowel or a consonant, and do so for nine
letters in total. They then compete by trying to form the longest possible word out of
those letters.

Validity of a word is decided by a dictionary, nothing else. There is no language processing.

All code is provided under the [MIT License](LICENSE).

## Approach

The input scramble, the letters on the board, are compared against a dictionary. The longest
matches are returned.

The basic approach is to convert words into 32-bit keys by setting a bit for each letter
that is in the word. You can imagine setting the first bit for the letter 'A', the third
bit for a 'C' and so on.

These keys, or bitmaps, are then used to quickly filter candidate matches, since any
dictionary candidate that contains a character that is NOT in the input scramble, can
be immediately skipped.

The 32-bit keys are arranged in tightly packed arrays, with each array containing keys only
for for a given word length. This makes it easy to rank the matches, and it's also quick to
look up the original characters in a parallel array -- which is required to validate words
that passed the initial bitmap filtering against false positives (e.g a word that used the
same letter more times than it occured in the input scramble).

## Performance

My machine can check 200,000 words in 0.4ms, which is fast enough that further optimization
doesn't seem warranted.

There is significant overhead in loading the dictionary, so using the command-line version to
query individual inputs will be much slower.

## Usage

### Command-Line Version

```bash
$ ./words-cli -file english-words-alpha.txt -limit 5 ibaselgsk
200005 words loaded from dictionary.
Looking for solutions to 'ibaselgsk', minlen=4, maxlen=9, limit=5, maxdist=-1
5 words found, 107320 rejected by falsebits, 153 rejected in validation, 0 rejected by distance. 107478 words checked. Duration: 179.257µs
0. 'kissable' dist=1
1. 'abseils' dist=2
2. 'algesis' dist=2
3. 'alsikes' dist=2
4. 'asslike' dist=2
```

### Service Version

```bash
$ ./words-server &>server.log &
$ curl -s localhost:8080/countdown/v1/words/LESRNOXIS | json_pp
```

```json
{
   "duration_ms" : 0.069932,
   "max_dist" : 4,
   "min_dist" : 1,
   "num_checked" : 42,
   "num_hits" : 5,
   "query" : "lesrnoxis",
   "words" : [
      "ironless",
      "lesions",
      "lesson",
      "snores",
      "roses"
   ]
}
```

## To Do

* Add TLS support
* Swagger API spec
* Example client to go with server
* Docker packaging?
