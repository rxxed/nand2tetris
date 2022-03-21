// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

@16384
D=A
@screenbegin
M=D

@24576
D=A
@screenend
M=D

(LOOP)
        // R0 = screenbegin
        @screenbegin
        D=M
        @R0
        M=D

        // If keyboard = 0, whiten screen
        @24576
        D=M
        @WHITEN
        D;JEQ

        // else, blacken screen
(BLACKEN)
        // if R0 == screenend, jump to LOOP
        @R0
        D=M
        @screenend
        D=M-D
        @LOOP
        D;JEQ

        // M[R0] = -1
        @R0
        A=M
        M=-1

        // R0++
        @R0
        D=M+1
        M=D

        @BLACKEN
        0;JMP

(WHITEN)
        // if R0 == screenend, jump to LOOP
        @R0
        D=M
        @screenend
        D=M-D
        @LOOP
        D;JEQ

        // M[R0] = 0
        @R0
        A=M
        M=0

        // R0++
        @R0
        D=M+1
        M=D

        @WHITEN
        0;JMP
