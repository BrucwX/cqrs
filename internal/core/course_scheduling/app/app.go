package app

import (
	// command handlers
	cabsence "cqrs/internal/core/course_scheduling/app/command/absence"

	cclassroom "cqrs/internal/core/course_scheduling/app/command/classroom"

	ccourse "cqrs/internal/core/course_scheduling/app/command/course"

	ccourseSlot "cqrs/internal/core/course_scheduling/app/command/courseSlot"

	ccourseSlotChange "cqrs/internal/core/course_scheduling/app/command/courseSlotChange"

	cenrollment "cqrs/internal/core/course_scheduling/app/command/enrollment"

	cqualification "cqrs/internal/core/course_scheduling/app/command/qualification"

	cstudent "cqrs/internal/core/course_scheduling/app/command/student"

	cteacher "cqrs/internal/core/course_scheduling/app/command/teacher"

	cmakeup "cqrs/internal/core/course_scheduling/app/command/makeup"

	// query handlers
	qabsence "cqrs/internal/core/course_scheduling/app/query/absence"

	qclassroom "cqrs/internal/core/course_scheduling/app/query/classroom"

	qcourse "cqrs/internal/core/course_scheduling/app/query/course"

	qcourseSlot "cqrs/internal/core/course_scheduling/app/query/courseSlot"

	qcourseSlotChange "cqrs/internal/core/course_scheduling/app/query/courseSlotChange"

	qenrollment "cqrs/internal/core/course_scheduling/app/query/enrollment"

	qmakeup "cqrs/internal/core/course_scheduling/app/query/makeup"

	qqualification "cqrs/internal/core/course_scheduling/app/query/qualification"

	qstudent "cqrs/internal/core/course_scheduling/app/query/student"

	qteacher "cqrs/internal/core/course_scheduling/app/query/teacher"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	Student       cstudent.Handler
	Teacher       cteacher.Handler
	Classroom     cclassroom.Handler
	Course        ccourse.Handler
	CourseSlot    ccourseSlot.Handler
	CourseSlotChg ccourseSlotChange.Handler
	Enrollment    cenrollment.Handler
	Absence       cabsence.Handler
	Makeup        cmakeup.Handler
	Qualification cqualification.Handler
}

type Queries struct {
	Student       qstudent.Handler
	Teacher       qteacher.Handler
	Classroom     qclassroom.Handler
	Course        qcourse.Handler
	CourseSlot    qcourseSlot.Handler
	CourseSlotChg qcourseSlotChange.Handler
	Enrollment    qenrollment.Handler
	Absence       qabsence.Handler
	Makeup        qmakeup.Handler
	Qualification qqualification.Handler
}
