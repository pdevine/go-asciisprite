package sprite

const Font_a = `
 X  
X X 
XXX 
X X 
X X `

const Font_b = `
XX  
X X 
XX  
X X 
XX  `

const Font_c = `    
 XX 
X   
X   
X   
 XX `

const Font_d = `
XX  
X X 
X X 
X X 
XXX `

const Font_e = `    
XXX 
X   
XX  
X   
XXX `

const Font_f = `
XXX 
X   
XX  
X   
X   `

const Font_g = `    
 XX 
X   
X   
X X 
XXX `

const Font_h = `    
X X 
X X 
XXX 
X X 
X X `

const Font_i = `    
XXX 
 X  
 X  
 X  
XXX `

const Font_j = `    
XXX 
 X  
 X  
 X  
XX  `

const Font_k = `    
X X 
X X 
XX  
X X 
X X `

const Font_l = `    
X   
X   
X   
X   
XXX `

const Font_m = `      
XXX 
XXX 
X X 
X X 
X X `

const Font_n = `      
XX  
X X 
X X 
X X 
X X `

const Font_o = `    
 XX 
X X 
X X 
X X 
XX  `

const Font_p = `    
XX  
X X 
XX  
X   
X   `

const Font_q = `    
 XX 
X X 
X X 
XXX 
  X `

const Font_r = `    
XX  
X X 
XX  
X X 
X X `

const Font_s = `    
 XX 
X   
XXX 
  X 
XX  `

const Font_t = `    
XXX 
 X  
 X  
 X  
 X  `

const Font_u = `    
X X 
X X 
X X 
X X 
 XX `

const Font_v = `    
X X 
X X 
X X 
X X 
 X  `

const Font_w = `    
X X 
X X 
X X 
XXX 
XXX `

const Font_x = `    
X X 
X X 
 X  
X X 
X X `

const Font_y = `    
X X 
X X 
XXX 
  X 
XXX `

const Font_z = `    
XXX 
  X 
 X  
X   
XXX `

const Font_0 = `    
 XX 
X X 
X X 
X X 
XX  `

const Font_1 = `    
 X  
XX  
 X  
 X  
XXX `

const Font_2 = `    
XX  
  X 
XXX 
X   
XXX `

const Font_3 = `    
XX  
  X 
XX  
  X 
XXX `

const Font_4 = `    
X X 
X X 
XXX 
  X 
  X `

const Font_5 = `    
XXX 
X   
XXX 
  X 
XX  `

const Font_6 = `    
 XX 
X   
XXX 
X X 
XX  `

const Font_7 = `    
XXX 
  X 
 X  
X   
X   `

const Font_8 = `    
 XX 
X X 
XXX 
X X 
XX  `

const Font_9 = `    
 XX 
X X 
XXX 
  X 
XXX `

const Font_period = `    
    
    
    
    
 X  `

const Font_comma = `    
    
    
    
 X  
 X  `

const Font_slash = `    
  X 
 X  
 X  
 X  
X   `

const Font_exclamation = `    
 X  
 X  
 X  
    
 X  `

const Font_dash = `    
    
    
XXX 
    
    `

// NewPakuFont provides a new font from based upon Paku Paku
func NewPakuFont() *Font {
	m := map[rune]string{
		'a': Font_a,
		'b': Font_b,
		'c': Font_c,
		'd': Font_d,
		'e': Font_e,
		'f': Font_f,
		'g': Font_g,
		'h': Font_h,
		'i': Font_i,
		'j': Font_j,
		'k': Font_k,
		'l': Font_l,
		'm': Font_m,
		'n': Font_n,
		'o': Font_o,
		'p': Font_p,
		'q': Font_q,
		'r': Font_r,
		's': Font_s,
		't': Font_t,
		'u': Font_u,
		'v': Font_v,
		'w': Font_w,
		'x': Font_x,
		'y': Font_y,
		'z': Font_z,
		'0': Font_0,
		'1': Font_1,
		'2': Font_2,
		'3': Font_3,
		'4': Font_4,
		'5': Font_5,
		'6': Font_6,
		'7': Font_7,
		'8': Font_8,
		'9': Font_9,
		'.': Font_period,
		',': Font_comma,
		'/': Font_slash,
		'!': Font_exclamation,
		'-': Font_dash,
	}

	return NewFont(m, 5, 7)
}

