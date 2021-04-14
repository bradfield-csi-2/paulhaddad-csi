default rel

section .text
global volume
volume:
	mulss xmm0, xmm0           ; r * 4
	mulss xmm0, xmm1           ; r**2 * h
	mulss xmm0, [pi_over_3]    ; pi/3 * r**2 * h
 	ret

section .rodata
pi_over_3: dd  1.0471975511965976
