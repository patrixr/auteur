release:
	git tag -a "v`cat ./VERSION`" -m "Release version `cat ./VERSION`"
	git push origin v`cat ./VERSION`
