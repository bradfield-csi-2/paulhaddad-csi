section .text
global fib
fib:
mov rax, rdi            ; Set return value to base case, n
cmp rdi, 1              ; are we at the base case of <= 1?
jle .base               ; if so, jump to 1 base case

push rbx                ; save rbx = n; save previous value
push r12                ; save r12 = result of call to fib(n-1) ; ask Oz about order here
mov rbx, rax            ; store n in callee saved register so we have access to it when procedure returns

dec rdi
call fib
mov r12, rax            ; store result of calling fib with n-1


lea rdi, [rbx-2]        ; calculate n-2
call fib
add rax, r12            ; add fib(n-1) + fib(n-2)

pop r12                 ; restore r12
pop rbx                 ; restore rbx

.base:
	ret                   ; return with 0 or 1
