// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 28, 306,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 3, 2, 3, 2,
	3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6,
	3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 8,
	3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10,
	3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 131, 10, 13, 3, 13, 3, 13, 3, 14,
	3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 5, 17, 142, 10, 17, 3, 17, 6,
	17, 145, 10, 17, 13, 17, 14, 17, 146, 3, 18, 3, 18, 3, 19, 3, 19, 3, 20,
	3, 20, 7, 20, 155, 10, 20, 12, 20, 14, 20, 158, 11, 20, 3, 21, 5, 21, 161,
	10, 21, 3, 21, 6, 21, 164, 10, 21, 13, 21, 14, 21, 165, 3, 21, 3, 21, 6,
	21, 170, 10, 21, 13, 21, 14, 21, 171, 3, 21, 3, 21, 3, 21, 6, 21, 177,
	10, 21, 13, 21, 14, 21, 178, 5, 21, 181, 10, 21, 5, 21, 183, 10, 21, 3,
	22, 5, 22, 186, 10, 22, 3, 22, 3, 22, 5, 22, 190, 10, 22, 3, 22, 6, 22,
	193, 10, 22, 13, 22, 14, 22, 194, 3, 22, 3, 22, 6, 22, 199, 10, 22, 13,
	22, 14, 22, 200, 3, 22, 3, 22, 3, 22, 3, 22, 3, 23, 5, 23, 208, 10, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 6, 23, 214, 10, 23, 13, 23, 14, 23, 215, 3,
	23, 3, 23, 3, 23, 3, 23, 6, 23, 222, 10, 23, 13, 23, 14, 23, 223, 3, 24,
	3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 6, 24, 235, 10,
	24, 13, 24, 14, 24, 236, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25,
	3, 25, 6, 25, 247, 10, 25, 13, 25, 14, 25, 248, 3, 26, 3, 26, 5, 26, 253,
	10, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 6, 28,
	263, 10, 28, 13, 28, 14, 28, 264, 3, 28, 3, 28, 3, 29, 3, 29, 6, 29, 271,
	10, 29, 13, 29, 14, 29, 272, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30,
	3, 31, 3, 31, 3, 31, 7, 31, 284, 10, 31, 12, 31, 14, 31, 287, 11, 31, 3,
	31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 32, 7, 32, 295, 10, 32, 12, 32, 14,
	32, 298, 11, 32, 3, 32, 3, 32, 3, 32, 5, 32, 303, 10, 32, 3, 32, 3, 32,
	2, 2, 33, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 19, 11,
	21, 12, 23, 13, 25, 14, 27, 15, 29, 16, 31, 17, 33, 18, 35, 2, 37, 2, 39,
	19, 41, 2, 43, 2, 45, 20, 47, 21, 49, 22, 51, 23, 53, 24, 55, 25, 57, 26,
	59, 2, 61, 27, 63, 28, 3, 2, 8, 8, 2, 38, 38, 44, 44, 48, 49, 66, 92, 97,
	97, 99, 124, 8, 2, 38, 38, 44, 44, 48, 59, 66, 92, 97, 97, 99, 124, 4,
	2, 45, 45, 47, 47, 4, 2, 50, 59, 99, 104, 3, 2, 36, 36, 4, 2, 12, 12, 15,
	15, 2, 327, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9,
	3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2,
	17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2,
	2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2,
	2, 2, 33, 3, 2, 2, 2, 2, 39, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3, 2,
	2, 2, 2, 49, 3, 2, 2, 2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2, 55, 3,
	2, 2, 2, 2, 57, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2, 2, 63, 3, 2, 2, 2, 3, 65,
	3, 2, 2, 2, 5, 72, 3, 2, 2, 2, 7, 77, 3, 2, 2, 2, 9, 82, 3, 2, 2, 2, 11,
	88, 3, 2, 2, 2, 13, 95, 3, 2, 2, 2, 15, 102, 3, 2, 2, 2, 17, 108, 3, 2,
	2, 2, 19, 114, 3, 2, 2, 2, 21, 118, 3, 2, 2, 2, 23, 122, 3, 2, 2, 2, 25,
	130, 3, 2, 2, 2, 27, 134, 3, 2, 2, 2, 29, 136, 3, 2, 2, 2, 31, 138, 3,
	2, 2, 2, 33, 141, 3, 2, 2, 2, 35, 148, 3, 2, 2, 2, 37, 150, 3, 2, 2, 2,
	39, 152, 3, 2, 2, 2, 41, 160, 3, 2, 2, 2, 43, 185, 3, 2, 2, 2, 45, 207,
	3, 2, 2, 2, 47, 225, 3, 2, 2, 2, 49, 238, 3, 2, 2, 2, 51, 252, 3, 2, 2,
	2, 53, 254, 3, 2, 2, 2, 55, 259, 3, 2, 2, 2, 57, 268, 3, 2, 2, 2, 59, 276,
	3, 2, 2, 2, 61, 280, 3, 2, 2, 2, 63, 290, 3, 2, 2, 2, 65, 66, 7, 111, 2,
	2, 66, 67, 7, 113, 2, 2, 67, 68, 7, 102, 2, 2, 68, 69, 7, 119, 2, 2, 69,
	70, 7, 110, 2, 2, 70, 71, 7, 103, 2, 2, 71, 4, 3, 2, 2, 2, 72, 73, 7, 118,
	2, 2, 73, 74, 7, 123, 2, 2, 74, 75, 7, 114, 2, 2, 75, 76, 7, 103, 2, 2,
	76, 6, 3, 2, 2, 2, 77, 78, 7, 104, 2, 2, 78, 79, 7, 119, 2, 2, 79, 80,
	7, 112, 2, 2, 80, 81, 7, 101, 2, 2, 81, 8, 3, 2, 2, 2, 82, 83, 7, 114,
	2, 2, 83, 84, 7, 99, 2, 2, 84, 85, 7, 116, 2, 2, 85, 86, 7, 99, 2, 2, 86,
	87, 7, 111, 2, 2, 87, 10, 3, 2, 2, 2, 88, 89, 7, 116, 2, 2, 89, 90, 7,
	103, 2, 2, 90, 91, 7, 117, 2, 2, 91, 92, 7, 119, 2, 2, 92, 93, 7, 110,
	2, 2, 93, 94, 7, 118, 2, 2, 94, 12, 3, 2, 2, 2, 95, 96, 7, 107, 2, 2, 96,
	97, 7, 111, 2, 2, 97, 98, 7, 114, 2, 2, 98, 99, 7, 113, 2, 2, 99, 100,
	7, 116, 2, 2, 100, 101, 7, 118, 2, 2, 101, 14, 3, 2, 2, 2, 102, 103, 7,
	110, 2, 2, 103, 104, 7, 113, 2, 2, 104, 105, 7, 101, 2, 2, 105, 106, 7,
	99, 2, 2, 106, 107, 7, 110, 2, 2, 107, 16, 3, 2, 2, 2, 108, 109, 7, 100,
	2, 2, 109, 110, 7, 110, 2, 2, 110, 111, 7, 113, 2, 2, 111, 112, 7, 101,
	2, 2, 112, 113, 7, 109, 2, 2, 113, 18, 3, 2, 2, 2, 114, 115, 7, 107, 2,
	2, 115, 116, 7, 53, 2, 2, 116, 117, 7, 52, 2, 2, 117, 20, 3, 2, 2, 2, 118,
	119, 7, 107, 2, 2, 119, 120, 7, 56, 2, 2, 120, 121, 7, 54, 2, 2, 121, 22,
	3, 2, 2, 2, 122, 123, 7, 104, 2, 2, 123, 124, 7, 56, 2, 2, 124, 125, 7,
	54, 2, 2, 125, 24, 3, 2, 2, 2, 126, 131, 7, 34, 2, 2, 127, 128, 7, 15,
	2, 2, 128, 131, 7, 12, 2, 2, 129, 131, 4, 11, 12, 2, 130, 126, 3, 2, 2,
	2, 130, 127, 3, 2, 2, 2, 130, 129, 3, 2, 2, 2, 131, 132, 3, 2, 2, 2, 132,
	133, 8, 13, 2, 2, 133, 26, 3, 2, 2, 2, 134, 135, 7, 42, 2, 2, 135, 28,
	3, 2, 2, 2, 136, 137, 7, 43, 2, 2, 137, 30, 3, 2, 2, 2, 138, 139, 7, 36,
	2, 2, 139, 32, 3, 2, 2, 2, 140, 142, 7, 47, 2, 2, 141, 140, 3, 2, 2, 2,
	141, 142, 3, 2, 2, 2, 142, 144, 3, 2, 2, 2, 143, 145, 4, 50, 59, 2, 144,
	143, 3, 2, 2, 2, 145, 146, 3, 2, 2, 2, 146, 144, 3, 2, 2, 2, 146, 147,
	3, 2, 2, 2, 147, 34, 3, 2, 2, 2, 148, 149, 9, 2, 2, 2, 149, 36, 3, 2, 2,
	2, 150, 151, 9, 3, 2, 2, 151, 38, 3, 2, 2, 2, 152, 156, 5, 35, 18, 2, 153,
	155, 5, 37, 19, 2, 154, 153, 3, 2, 2, 2, 155, 158, 3, 2, 2, 2, 156, 154,
	3, 2, 2, 2, 156, 157, 3, 2, 2, 2, 157, 40, 3, 2, 2, 2, 158, 156, 3, 2,
	2, 2, 159, 161, 7, 47, 2, 2, 160, 159, 3, 2, 2, 2, 160, 161, 3, 2, 2, 2,
	161, 163, 3, 2, 2, 2, 162, 164, 4, 50, 59, 2, 163, 162, 3, 2, 2, 2, 164,
	165, 3, 2, 2, 2, 165, 163, 3, 2, 2, 2, 165, 166, 3, 2, 2, 2, 166, 182,
	3, 2, 2, 2, 167, 169, 7, 48, 2, 2, 168, 170, 4, 50, 59, 2, 169, 168, 3,
	2, 2, 2, 170, 171, 3, 2, 2, 2, 171, 169, 3, 2, 2, 2, 171, 172, 3, 2, 2,
	2, 172, 180, 3, 2, 2, 2, 173, 174, 7, 103, 2, 2, 174, 176, 9, 4, 2, 2,
	175, 177, 4, 50, 59, 2, 176, 175, 3, 2, 2, 2, 177, 178, 3, 2, 2, 2, 178,
	176, 3, 2, 2, 2, 178, 179, 3, 2, 2, 2, 179, 181, 3, 2, 2, 2, 180, 173,
	3, 2, 2, 2, 180, 181, 3, 2, 2, 2, 181, 183, 3, 2, 2, 2, 182, 167, 3, 2,
	2, 2, 182, 183, 3, 2, 2, 2, 183, 42, 3, 2, 2, 2, 184, 186, 7, 47, 2, 2,
	185, 184, 3, 2, 2, 2, 185, 186, 3, 2, 2, 2, 186, 189, 3, 2, 2, 2, 187,
	188, 7, 50, 2, 2, 188, 190, 7, 122, 2, 2, 189, 187, 3, 2, 2, 2, 189, 190,
	3, 2, 2, 2, 190, 192, 3, 2, 2, 2, 191, 193, 9, 5, 2, 2, 192, 191, 3, 2,
	2, 2, 193, 194, 3, 2, 2, 2, 194, 192, 3, 2, 2, 2, 194, 195, 3, 2, 2, 2,
	195, 196, 3, 2, 2, 2, 196, 198, 7, 48, 2, 2, 197, 199, 9, 5, 2, 2, 198,
	197, 3, 2, 2, 2, 199, 200, 3, 2, 2, 2, 200, 198, 3, 2, 2, 2, 200, 201,
	3, 2, 2, 2, 201, 202, 3, 2, 2, 2, 202, 203, 7, 114, 2, 2, 203, 204, 9,
	4, 2, 2, 204, 205, 4, 50, 59, 2, 205, 44, 3, 2, 2, 2, 206, 208, 7, 47,
	2, 2, 207, 206, 3, 2, 2, 2, 207, 208, 3, 2, 2, 2, 208, 209, 3, 2, 2, 2,
	209, 210, 7, 50, 2, 2, 210, 211, 7, 122, 2, 2, 211, 213, 3, 2, 2, 2, 212,
	214, 4, 50, 59, 2, 213, 212, 3, 2, 2, 2, 214, 215, 3, 2, 2, 2, 215, 213,
	3, 2, 2, 2, 215, 216, 3, 2, 2, 2, 216, 217, 3, 2, 2, 2, 217, 218, 7, 114,
	2, 2, 218, 219, 7, 45, 2, 2, 219, 221, 3, 2, 2, 2, 220, 222, 4, 50, 59,
	2, 221, 220, 3, 2, 2, 2, 222, 223, 3, 2, 2, 2, 223, 221, 3, 2, 2, 2, 223,
	224, 3, 2, 2, 2, 224, 46, 3, 2, 2, 2, 225, 226, 7, 113, 2, 2, 226, 227,
	7, 104, 2, 2, 227, 228, 7, 104, 2, 2, 228, 229, 7, 117, 2, 2, 229, 230,
	7, 103, 2, 2, 230, 231, 7, 118, 2, 2, 231, 232, 7, 63, 2, 2, 232, 234,
	3, 2, 2, 2, 233, 235, 4, 50, 59, 2, 234, 233, 3, 2, 2, 2, 235, 236, 3,
	2, 2, 2, 236, 234, 3, 2, 2, 2, 236, 237, 3, 2, 2, 2, 237, 48, 3, 2, 2,
	2, 238, 239, 7, 99, 2, 2, 239, 240, 7, 110, 2, 2, 240, 241, 7, 107, 2,
	2, 241, 242, 7, 105, 2, 2, 242, 243, 7, 112, 2, 2, 243, 244, 7, 63, 2,
	2, 244, 246, 3, 2, 2, 2, 245, 247, 4, 50, 59, 2, 246, 245, 3, 2, 2, 2,
	247, 248, 3, 2, 2, 2, 248, 246, 3, 2, 2, 2, 248, 249, 3, 2, 2, 2, 249,
	50, 3, 2, 2, 2, 250, 253, 5, 41, 21, 2, 251, 253, 5, 43, 22, 2, 252, 250,
	3, 2, 2, 2, 252, 251, 3, 2, 2, 2, 253, 52, 3, 2, 2, 2, 254, 255, 7, 61,
	2, 2, 255, 256, 7, 63, 2, 2, 256, 257, 5, 41, 21, 2, 257, 258, 7, 61, 2,
	2, 258, 54, 3, 2, 2, 2, 259, 260, 7, 61, 2, 2, 260, 262, 7, 66, 2, 2, 261,
	263, 4, 50, 59, 2, 262, 261, 3, 2, 2, 2, 263, 264, 3, 2, 2, 2, 264, 262,
	3, 2, 2, 2, 264, 265, 3, 2, 2, 2, 265, 266, 3, 2, 2, 2, 266, 267, 7, 61,
	2, 2, 267, 56, 3, 2, 2, 2, 268, 270, 7, 61, 2, 2, 269, 271, 4, 50, 59,
	2, 270, 269, 3, 2, 2, 2, 271, 272, 3, 2, 2, 2, 272, 270, 3, 2, 2, 2, 272,
	273, 3, 2, 2, 2, 273, 274, 3, 2, 2, 2, 274, 275, 7, 61, 2, 2, 275, 58,
	3, 2, 2, 2, 276, 277, 7, 94, 2, 2, 277, 278, 9, 5, 2, 2, 278, 279, 9, 5,
	2, 2, 279, 60, 3, 2, 2, 2, 280, 285, 7, 36, 2, 2, 281, 284, 5, 59, 30,
	2, 282, 284, 10, 6, 2, 2, 283, 281, 3, 2, 2, 2, 283, 282, 3, 2, 2, 2, 284,
	287, 3, 2, 2, 2, 285, 283, 3, 2, 2, 2, 285, 286, 3, 2, 2, 2, 286, 288,
	3, 2, 2, 2, 287, 285, 3, 2, 2, 2, 288, 289, 7, 36, 2, 2, 289, 62, 3, 2,
	2, 2, 290, 291, 7, 61, 2, 2, 291, 292, 7, 61, 2, 2, 292, 296, 3, 2, 2,
	2, 293, 295, 10, 7, 2, 2, 294, 293, 3, 2, 2, 2, 295, 298, 3, 2, 2, 2, 296,
	294, 3, 2, 2, 2, 296, 297, 3, 2, 2, 2, 297, 302, 3, 2, 2, 2, 298, 296,
	3, 2, 2, 2, 299, 300, 7, 15, 2, 2, 300, 303, 7, 12, 2, 2, 301, 303, 7,
	12, 2, 2, 302, 299, 3, 2, 2, 2, 302, 301, 3, 2, 2, 2, 303, 304, 3, 2, 2,
	2, 304, 305, 8, 32, 2, 2, 305, 64, 3, 2, 2, 2, 29, 2, 130, 141, 146, 156,
	160, 165, 171, 178, 180, 182, 185, 189, 194, 200, 207, 215, 223, 236, 248,
	252, 264, 272, 283, 285, 296, 302, 3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'module'", "'type'", "'func'", "'param'", "'result'", "'import'",
	"'local'", "'block'", "'i32'", "'i64'", "'f64'", "", "'('", "')'", "'\"'",
}

var lexerSymbolicNames = []string{
	"", "ModuleWord", "TypeWord", "FuncWord", "ParamWord", "ResultWord", "ImportWord",
	"LocalWord", "BlockWord", "I32", "I64", "F64", "Whitespace", "Lparen",
	"Rparen", "Quote", "Num", "Ident", "HexPointer", "Offset", "Align", "ConstValue",
	"ConstAnnotation", "BlockAnnotation", "TypeAnnotation", "QuotedString",
	"Comment",
}

var lexerRuleNames = []string{
	"ModuleWord", "TypeWord", "FuncWord", "ParamWord", "ResultWord", "ImportWord",
	"LocalWord", "BlockWord", "I32", "I64", "F64", "Whitespace", "Lparen",
	"Rparen", "Quote", "Num", "IdentFirst", "IdentAfter", "Ident", "IntConst",
	"FloatConst", "HexPointer", "Offset", "Align", "ConstValue", "ConstAnnotation",
	"BlockAnnotation", "TypeAnnotation", "HexByteValue", "QuotedString", "Comment",
}

type jsstripLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewjsstripLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *jsstripLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewjsstripLexer(input antlr.CharStream) *jsstripLexer {
	l := new(jsstripLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "jsstrip.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// jsstripLexer tokens.
const (
	jsstripLexerModuleWord      = 1
	jsstripLexerTypeWord        = 2
	jsstripLexerFuncWord        = 3
	jsstripLexerParamWord       = 4
	jsstripLexerResultWord      = 5
	jsstripLexerImportWord      = 6
	jsstripLexerLocalWord       = 7
	jsstripLexerBlockWord       = 8
	jsstripLexerI32             = 9
	jsstripLexerI64             = 10
	jsstripLexerF64             = 11
	jsstripLexerWhitespace      = 12
	jsstripLexerLparen          = 13
	jsstripLexerRparen          = 14
	jsstripLexerQuote           = 15
	jsstripLexerNum             = 16
	jsstripLexerIdent           = 17
	jsstripLexerHexPointer      = 18
	jsstripLexerOffset          = 19
	jsstripLexerAlign           = 20
	jsstripLexerConstValue      = 21
	jsstripLexerConstAnnotation = 22
	jsstripLexerBlockAnnotation = 23
	jsstripLexerTypeAnnotation  = 24
	jsstripLexerQuotedString    = 25
	jsstripLexerComment         = 26
)
