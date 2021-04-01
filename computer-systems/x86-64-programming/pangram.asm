section .text
global pangram
pangram:
	xor rax, rax
	xor r8, r8                 ; string increment

.loop:
	movzx r9, byte [rdi+r8]   ; get character
	cmp r9, 0                 ; are we at the end of the string?
	je .check
	inc r8
	cmp r9, 0x41              ; skip letters less than 'A'
	jl .loop
	or r9, 0x20               ; set fifth bit to force lowercase
	sub r9, 'a'               ; get int from letter
	bts rax, r9               ; shift bit
	jmp .loop

.check:
cmp rax, 0x3ffffff          ; are all the letters shifted: 0b11111...
je .done
xor rax, rax                ; if not, return 0

.done:
	ret
