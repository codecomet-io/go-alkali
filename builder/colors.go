package builder

// Stuff to be put back into isovaline commander

import (
	"fmt"
	"os"
)

func Colors() string {
	solarDark()

	fromEnv := os.Getenv("CODECOMET_COLORS")
	if fromEnv == "" {
		fromEnv = solarToBk()
	}

	// https://github.com/moby/buildkit/issues/3537
	nocol := os.Getenv("NO_COLOR")
	if nocol == "" {
		os.Unsetenv("NO_COLOR")
	}

	/*
		const n = 20
		builder := aec.EmptyBuilder

		up2 := aec.Up(2)
		col := aec.Column(n + 2)
		bar := aec.Color8BitF(aec.NewRGB8Bit(64, 255, 64))
		label := builder.CyanB().LightRedF().Underline().With(col).Right(1).ANSI

		// for up2
		fmt.Println()
		fmt.Println()

		for i := 0; i <= n; i++ {
			fmt.Print(up2)
			fmt.Println(label.Apply(fmt.Sprint(i, "/", n)))
			fmt.Print("[")
			fmt.Print(bar.Apply(strings.Repeat("=", i)))
			fmt.Println(col.Apply("]"))
			time.Sleep(100 * time.Millisecond)
		}

	*/ //nolint:dupword

	return fromEnv
}

func solarToBk() string {
	return fmt.Sprintf("run=%s:warning=%s:error=%s:cancel=%s", SolViolet, SolOrange, SolRed, SolMagenta)
}

const (
	SolBase03  = "0,43,54"
	SolBase02  = "7,54,66"
	SolBase01  = "88,110,117"
	SolBase00  = "101,123,131"
	SolBase0   = "131,148,150"
	SolBase1   = "147,161,161"
	SolBase2   = "238,232,213"
	SolBase3   = "253,246,227"
	SolYellow  = "181,137,0"
	SolOrange  = "203,75,22"
	SolRed     = "220,50,47"
	SolMagenta = "211,54,130"
	SolViolet  = "108,113,196"
	SolBlue    = "38,139,210"
	SolCyan    = "42,161,152"
	SolGreen   = "133,153,0"
)

var (
	SolBodyText             string //nolint:gochecknoglobals
	SolEmphasize            string //nolint:gochecknoglobals
	SolComments             string //nolint:gochecknoglobals
	SolBackgroundHighlights string //nolint:gochecknoglobals
	SolBackground           string //nolint:gochecknoglobals
	SolError                string //nolint:gochecknoglobals
	SolWarning              string //nolint:gochecknoglobals
	SolInfo                 string //nolint:gochecknoglobals
	SolDebug                string //nolint:gochecknoglobals
)

func solarDark() {
	SolBodyText = SolBase0
	SolEmphasize = SolBase1
	SolComments = SolBase01
	SolBackgroundHighlights = SolBase02
	SolBackground = SolBase03
	SolError = SolRed
	SolWarning = SolOrange
	SolInfo = SolGreen
	SolDebug = SolBodyText
}

/*
func solarLight() {
	SolBodyText = SolBase00
	SolEmphasize = SolBase01
	SolComments = SolBase1
	SolBackgroundHighlights = SolBase2
	SolBackground = SolBase3
	SolError = SolRed
	SolWarning = SolOrange
	SolInfo = SolGreen
	SolDebug = SolBodyText
}

*/

/*

#SOLARIZED HEX     16/8 TERMCOL  XTERM/HEX   L*A*B      RGB         HSB
#--------- ------- ---- -------  ----------- ---------- ----------- -----------
#base03    #002b36  8/4 brblack  234 #1c1c1c 15 -12 -12     0  43  54 193 100  21
#base02    #073642  0/4 black    235 #262626 20 -12 -12     7  54  66 192  90  26
#base01    #586e75 10/7 brgreen  240 #585858 45 -07 -07    88 110 117 194  25  46
#base00    #657b83 11/7 bryellow 241 #626262 50 -07 -07   101 123 131 195  23  51
#base0     #839496 12/6 brblue   244 #808080 60 -06 -03   131 148 150 186  13  59
#base1     #93a1a1 14/4 brcyan   245 #8a8a8a 65 -05 -02   147 161 161 180   9  63
#base2     #eee8d5  7/7 white    254 #e4e4e4 92 -00  10   238 232 213  44  11  93
#base3     #fdf6e3 15/7 brwhite  230 #ffffd7 97  00  10   253 246 227  44  10  99

#yellow    #b58900  3/3 yellow   136 #af8700 60  10  65   181 137   0  45 100  71
#orange    #cb4b16  9/3 brred    166 #d75f00 50  50  55   203  75  22  18  89  80
#red       #dc322f  1/1 red      160 #d70000 50  65  45   220  50  47   1  79  86
#magenta   #d33682  5/5 magenta  125 #af005f 50  65 -05   211  54 130 331  74  83
#violet    #6c71c4 13/5 brmagenta 61 #5f5faf 50  15 -45   108 113 196 237  45  77
#blue      #268bd2  4/4 blue      33 #0087ff 55 -10 -45    38 139 210 205  82  82
#cyan      #2aa198  6/6 cyan      37 #00afaf 60 -35 -05    42 161 152 175  74  63
#green     #859900  2/2 green     64 #5f8700 60 -20  65   133 153   0  68 100  60
*/
